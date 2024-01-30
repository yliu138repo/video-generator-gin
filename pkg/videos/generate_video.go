package videos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yliu138repo/video-generator-gin/pkg/common/system"
)

type GenerateVideoBody struct {
	BgmMusic     string   `json:"bgmMusic"`
	CoverPage    string   `json:"coverPage"`
	VideoSrcList []string `json:"videoSrcList"`
	Title        string   `json:"title"`
}

// Validate file path format and if file exists from the path
func checkInput(body GenerateVideoBody) (string, error) {
	bgmErr := CheckPathExist(body.BgmMusic)
	if bgmErr != nil {
		return body.BgmMusic, bgmErr
	}

	coverErr := CheckPathExist(body.CoverPage)
	if coverErr != nil {
		return body.CoverPage, coverErr
	}

	if len(body.VideoSrcList) == 0 {
		errMsg := "no video src provided"
		return errMsg, errors.New(errMsg)
	}

	for _, videoSrc := range body.VideoSrcList {
		videoDirErr := CheckPathExist(videoSrc)
		if videoDirErr != nil {
			return videoSrc, videoDirErr
		}
	}

	return "", nil
}

// @BasePath /api/v1
// A POST function which generates video based on video sources and themes (background, cover and music) selected. All params should be absolute path
// @Summary Accept user-provided videos, images and themes, and generate video for user to download
// @Schemes
// @Description A POST function which generates video based on video sources and themes (background, cover and music) selected. All params should be absolute path
// @Tags video
// @Accept json
// @Produce json
// @Param req body videos.GenerateVideoBody true "GenerateVideoBody"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} videos.ErrorMessage
// @Router /videos [POST]
func (h handler) GenerateVideo(c *gin.Context) {
	body := GenerateVideoBody{}

	// Get requests's body
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	filePath, err := checkInput(body)
	if err != nil {
		log.Printf("error: %s: %+v", filePath, err)
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: filePath + ": " + err.Error(),
		})
		return
	}

	outputPath, pid, videoErr := GenerateVideo(c.Request.Context(), body)
	if videoErr != nil {
		log.Printf("video generating error: %+v", videoErr)
		c.JSON(http.StatusInternalServerError, ErrorMessage{
			ErrorMsg: fmt.Sprintf("video generating error: %+v", videoErr.Error()),
		})
		return
	}

	publicIP := GetIP2()
	c.JSON(http.StatusOK, map[string]interface{}{
		"outputPath": fmt.Sprintf("http://%s%s/videos/download?file_path=%s", publicIP, viper.Get("PORT").(string), url.QueryEscape(outputPath)),
		"pid":        pid,
		"ip":         publicIP,
	})
}

type ProcessResult struct {
	ErrorCode      int   `json:"errorCode"`
	Error          error `json:"error"`
	ProcessSucceed bool  `json:"processSucceed"`
}

// Genreate a new video based on the input
// Note it will remove the file first
func GenerateVideo(ctx context.Context, body GenerateVideoBody) (string, int, error) {
	outputPath := filepath.Join(filepath.Dir(body.VideoSrcList[0]), "output.mp4")
	// Remove if exists
	rmErr := system.RemoveFileIfExists(outputPath)
	if rmErr != nil {
		return outputPath, -1, rmErr
	}

	// For video concatenation adjustment
	aspectRatioWidth, aspectRatioHeight, sar := 1280, 720, 1
	videoAdjustCmd := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,setsar=%d", aspectRatioWidth, aspectRatioHeight, aspectRatioWidth, aspectRatioHeight, sar)

	srcLen := len(body.VideoSrcList)
	inputCmd, filterComplex := "", ""
	for i, src := range body.VideoSrcList {
		if checkImage(src) {
			inputCmd = inputCmd + " -loop 1 -framerate 24 -t 10 -i " + src
		} else {
			inputCmd = inputCmd + " -i " + src
		}

		filterComplex = filterComplex + fmt.Sprintf("[%d:v]%s[vid%d];", i, videoAdjustCmd, i)
	}

	// Keep add cover and bgm
	inputCmd = inputCmd + fmt.Sprintf(" -i %s -i %s", body.CoverPage, body.BgmMusic)
	// add cover
	filterComplex = filterComplex + fmt.Sprintf("[%d:v]%s[cover];[cover]", srcLen, videoAdjustCmd)
	// concatenate the video streams
	for i := 0; i < srcLen; i++ {
		filterComplex = filterComplex + fmt.Sprintf("[vid%d]", i)
	}
	filterComplexVideo := filterComplex + fmt.Sprintf("concat=n=%d:v=1:a=0[v]", (srcLen+1))
	// append audio stream
	filterComplexAudio := fmt.Sprintf("[%d:a]amerge=inputs=1[a]", (srcLen + 1))

	// To compose the arguments of ffmpeg
	framerate := viper.Get("FRAME_RATE").(string)
	args := "-y -framerate " + framerate + " -pix_fmt yuv420p "
	args = args + inputCmd + " -filter_complex " + filterComplexVideo + " -filter_complex " + filterComplexAudio + ` -map [v] -map [a] -ac 2 -shortest ` + outputPath

	fmt.Printf("%s &&&\n", args)
	argsAr := strings.Fields(args)

	pid, err := RunCommand("ffmpeg", argsAr, func(cmd *exec.Cmd, cmdErr error) {
		log.Printf("Inserting record - PID: %d, Process succeed: %+v,  error code: %+v, error: %+v\n", cmd.Process.Pid, cmd.ProcessState.Success(), cmd.ProcessState.ExitCode(), cmdErr)
		currentWD, err := system.CurrentWD()
		if err != nil {
			log.Printf("Failed to get current WD: %+v\n", err)
		}
		resultFilePath := fmt.Sprintf("%s/result.json", currentWD)
		resultJson, err := system.ReadJson[ProcessResult](resultFilePath)
		pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
		if err != nil {
			resultMap := map[string]interface{}{}
			resultMap[pidStr] = ProcessResult{
				ErrorCode:      cmd.ProcessState.ExitCode(),
				Error:          cmdErr,
				ProcessSucceed: cmd.ProcessState.Success(),
			}
			writeErr := system.WriteJson(resultFilePath, resultMap)
			if writeErr != nil {
				log.Printf("Failed to write json file to %s: %+v\n", resultFilePath, writeErr)
			}
		} else {
			resultJson[pidStr] = ProcessResult{
				ErrorCode:      cmd.ProcessState.ExitCode(),
				Error:          cmdErr,
				ProcessSucceed: cmd.ProcessState.Success(),
			}
			writeErr := system.WriteJson(resultFilePath, resultJson)
			if writeErr != nil {
				log.Printf("Failed to write json file to %s: %+v\n", resultFilePath, writeErr)
			}
		}

	})
	return outputPath, pid, err
}

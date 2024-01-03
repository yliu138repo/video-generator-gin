package videos

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
		errMsg := "No video src provided"
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

// A POST function which generates video based on
// All params should be absolute path
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

	outputPath, videoErr := GenerateVideo(body)
	if videoErr != nil {
		log.Printf("video generating error: %+v", videoErr)
		c.JSON(http.StatusInternalServerError, ErrorMessage{
			ErrorMsg: fmt.Sprintf("video generating error: %+v", videoErr.Error()),
		})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(outputPath))
	c.Header("Content-Type", "application/octet-stream")
	c.File(outputPath)
	// c.JSON(http.StatusOK, "Video was generated successfully")
}

// Genreate a new video based on the input
func GenerateVideo(body GenerateVideoBody) (string, error) {
	outputPath := filepath.Join(filepath.Dir(body.VideoSrcList[0]), "output.mp4")

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
	filterComplex = filterComplex + fmt.Sprintf("concat=n=%d:v=1:a=0[v];", (srcLen+1))
	// append audio stream
	filterComplex = filterComplex + fmt.Sprintf("[%d:a]amerge=inputs=1[a]", (srcLen+1))

	// To compose the arguments of ffmpeg
	framerate := viper.Get("FRAME_RATE").(string)
	args := "-y -framerate " + framerate + " -pix_fmt yuv420p "
	args = args + inputCmd + " -filter_complex " + filterComplex + ` -map [v] -map [a] -ac 2 -shortest ` + outputPath

	argsAr := strings.Fields(args)

	err := RunCommand("ffmpeg", argsAr)
	return outputPath, err
}

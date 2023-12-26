package videos

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type GenerateVideoBody struct {
	GgmMusic  string `json:"bgmMusic"`
	CoverPage string `json:"coverPage"`
	VideoDir  string `json:"videoDir"`
}

func checkPathExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file path %s does not exist", path)
		}
		return err
	}
	return nil
}

// Validate file path format and if file exists from the path
func checkInput(body GenerateVideoBody) (string, error) {
	bgmErr := checkPathExist(body.GgmMusic)
	if bgmErr != nil {
		return body.GgmMusic, bgmErr
	}

	coverErr := checkPathExist(body.CoverPage)
	if coverErr != nil {
		return body.CoverPage, coverErr
	}

	videoDirErr := checkPathExist(body.VideoDir)
	if videoDirErr != nil {
		return body.VideoDir, videoDirErr
	}

	return "", nil
}

type ErrorMessage struct {
	ErrorMsg string `json:"errorMessage"`
}

// Generate video based
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
	outputPath := filepath.Join(body.VideoDir, "out.mp4")
	framerate := viper.Get("FRAME_RATE").(string)
	args := []string{"-y", "-framerate", framerate, "-i", filepath.Join(body.VideoDir, "%d.jpg"), "-i", body.GgmMusic, "-c:v", "libx264", "-pix_fmt", "yuv420p", "-vf", "scale=320:240", "-t", "15", "-shortest", outputPath}
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return outputPath, err
	}

	err = cmd.Start()
	if err != nil {
		return outputPath, err
	}

	bytes, err := io.ReadAll(stdOut)
	log.Println(bytes)
	if err != nil {
		return outputPath, err
	}

	log.Printf("Waiting for command to finish...")
	log.Printf("Process id is %v", cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("Exit error is %+v, error code: %v\n", exitError, exitError.ExitCode())
			return outputPath, exitError
		}
	}
	return outputPath, nil
}

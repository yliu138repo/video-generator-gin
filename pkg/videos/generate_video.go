package videos

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hellokvn/go-gin-api-medium/pkg/common/system"
	"github.com/spf13/viper"
)

type GenerateVideoBody struct {
	GgmMusic  string `json:"bgmMusic"`
	CoverPage string `json:"coverPage"`
	VideoDir  string `json:"videoDir"`
	Title     string `json:"title"`
}

// Validate file path format and if file exists from the path
func checkInput(body GenerateVideoBody) (string, error) {
	bgmErr := CheckPathExist(body.GgmMusic)
	if bgmErr != nil {
		return body.GgmMusic, bgmErr
	}

	coverErr := CheckPathExist(body.CoverPage)
	if coverErr != nil {
		return body.CoverPage, coverErr
	}

	videoDirErr := CheckPathExist(body.VideoDir)
	if videoDirErr != nil {
		return body.VideoDir, videoDirErr
	}

	return "", nil
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
	cpErr := system.CopyFile(body.CoverPage, filepath.Join(body.VideoDir, "0.jpg"))
	if cpErr != nil {
		log.Printf("Failed to copy file from %s to %s\n", body.CoverPage, filepath.Join(body.VideoDir, "0.jpg"))
		return outputPath, cpErr
	}

	err := RunCommand("ffmpeg", args)
	return outputPath, err
}

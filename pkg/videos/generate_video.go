package videos

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type GenerateVideoBody struct {
	GgmMusic   string   `json:"bgmMusic"`
	CoverPage  string   `json:"coverPage"`
	VideoPaths []string `json:"videoPaths"`
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

	for _, path := range body.VideoPaths {
		videoErr := checkPathExist(path)
		if videoErr != nil {
			return path, videoErr
		}
	}

	return "", nil
}

type ErrorMessage struct {
	ErrorMsg string `json:"errorMessage"`
}

// Generate video based
// All params should be absolute path
func (h handler) GenerateVideo(c *gin.Context) {
	outputPath := "D:/videos/exported.mp4"
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

	c.JSON(http.StatusOK, body)
}

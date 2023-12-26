package videos

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type GenerateCoverPageBody struct {
	CoverPage string `json:"coverPage"`
	Title     string `json:"title"`
}

func checkCoverPageInput(body GenerateCoverPageBody) error {
	coverErr := CheckPathExist(body.CoverPage)
	if coverErr != nil {
		return coverErr
	}

	if body.Title == "" {
		return errors.New("tilte is invalid")
	}

	return nil
}

func (h handler) GenerateCoverPage(c *gin.Context) {
	body := GenerateCoverPageBody{}
	// Get requests's body
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := checkCoverPageInput(body)
	if err != nil {
		log.Printf("error: %+v", err)
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: err.Error(),
		})
		return
	}

	outputPath, coverErr := GenerateCoverPage(body, "white", 100, 1002, 100)
	if coverErr != nil {
		log.Printf("video generating error: %+v", coverErr)
		c.JSON(http.StatusInternalServerError, ErrorMessage{
			ErrorMsg: fmt.Sprintf("video generating error: %+v", coverErr.Error()),
		})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(outputPath))
	c.Header("Content-Type", "application/octet-stream")
	c.File(outputPath)
}

// Generate cover page with tile and styles.
func GenerateCoverPage(body GenerateCoverPageBody, fontColor string, fontSize int64, x int64, y int64) (string, error) {
	outputPath := filepath.Join(filepath.Dir(body.CoverPage), "cover-modified.jpg")
	args := []string{"-y", "-i", body.CoverPage, "-vf", fmt.Sprintf("drawtext=text='%s':fontcolor=%s:fontsize=%d:x=%d:y=%d:", body.Title, fontColor, fontSize, x, y), outputPath}

	err := RunCommand("ffmpeg", args)
	return outputPath, err
}

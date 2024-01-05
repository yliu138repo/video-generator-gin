package videos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type GenerateCoverPageBody struct {
	CoverPage     string `json:"coverPage"`
	Title         string `json:"title"`
	DestPath      string `json:"destPath"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	FadeInDuation string `json:"fadeInDuration"`
	FontColor     string `json:"fontColor"`
	FonSize       string `json:"fontSize"`
	X             string `json:"x"`
	Y             string `json:"y"`
}

func checkCoverInput(body GenerateCoverPageBody) error {
	coverErr := CheckPathExist(body.CoverPage)
	if coverErr != nil {
		return coverErr
	}

	if body.Title == "" {
		return errors.New("tilte is invalid")
	}

	return nil
}

// @BasePath /api/v1
// A POST function which generates cover videos based on user input, e.g. font, size, styles etc.
// @Summary Generates cover videos based on user input, e.g. font, size, styles etc.
// @Schemes
// @Description A POST function which generates cover videos based on user input, e.g. font, size, styles etc.
// @Tags video
// @Accept json
// @Produce json
// @Param req body videos.GenerateCoverPageBody true "GenerateCoverPageBody"
// @Success 200 {string} video file content
// @Failure 400 {string} media file does not exist  "Bad requests"
// @Router /videos/cover [POST]
func (h handler) GenerateCoverPage(c *gin.Context) {
	body := GenerateCoverPageBody{}
	// Get requests's body
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := checkCoverInput(body)
	if err != nil {
		log.Printf("error: %+v", err)
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: err.Error(),
		})
		return
	}

	outputPath, coverErr := GenerateCoverVideo(c.Request.Context(), body)
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

// Generate cover page with tile and styles with fadein effects
// when x and y is empty string, will put the title in the middle of the screen
func GenerateCoverVideo(ctx context.Context, body GenerateCoverPageBody) (string, error) {
	outputPath := filepath.Join(body.DestPath)
	if body.X == "" {
		body.X = "(w-text_w)/2"
	}

	if body.Y == "" {
		body.Y = "(h-text_h)/2"
	}

	args := []string{
		"-y", "-i", body.CoverPage, "-vf",
		fmt.Sprintf("drawtext=text='%s':fontcolor=%s:fontsize=%s:x=%s:y=%s:enable='between(t,%s,%s)',fade=t=in:st=%s:d=%s:alpha=1", body.Title, body.FontColor, body.FonSize, body.X, body.Y, body.StartTime, body.EndTime, body.StartTime, body.FadeInDuation),
		"-codec:a", "copy",
		outputPath,
	}

	err := RunCommandContext(ctx, "ffmpeg", args)
	return outputPath, err
}

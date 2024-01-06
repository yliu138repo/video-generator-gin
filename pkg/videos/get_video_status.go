package videos

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yliu138repo/video-generator-gin/pkg/common/system"
)

// @BasePath /api/v1
// A GET function which fetch the video generating status
// @Summary Accespt query params PID: process ID, ip: IP address of the server, outputPath: the pre-provided output path for the video
// @Schemes
// @Description A GET function which fetch the video generating status, Output path should be absolute.
// @Tags video
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /videos [GET]
func (h handler) GetVideoStatus(c *gin.Context) {
	// Validate input
	_, ok := c.GetQuery("ip")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: "ip not provided",
		})
		return
	}

	outputPath, ok := c.GetQuery("outputPath")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: "outputPath not provided",
		})
		return
	}

	pidStr, ok := c.GetQuery("pid")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: "pid not provided",
		})
		return
	}

	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: "pid is not a valid format",
		})
		return
	}

	if system.FileExists(outputPath) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"videoGenerated": true,
			"processStatus":  "Done",
		})
		return
	} else {
		_, err := os.FindProcess(int(pid))
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorMessage{
				ErrorMsg: fmt.Sprintf("Failed to run the PID: %+v", err),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"videoGenerated": false,
				"processStatus":  "Running",
			})
		}
	}
}

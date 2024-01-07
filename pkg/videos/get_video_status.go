package videos

import (
	"fmt"
	"log"
	"net/http"

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
// @Failure 400 {object} videos.ErrorMessage
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
	fmt.Println(outputPath)

	pidStr, ok := c.GetQuery("pid")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: "pid not provided",
		})
		return
	}

	// pid, err := strconv.ParseInt(pidStr, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, ErrorMessage{
	// 		ErrorMsg: "pid is not a valid format",
	// 	})
	// 	return
	// }

	currentWD, err := system.CurrentWD()
	if err != nil {
		log.Printf("Failed to get current WD: %+v\n", err)
	}
	resultFilePath := fmt.Sprintf("%s/result.json", currentWD)
	resultJson, err := system.ReadJson[ProcessResult](resultFilePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{
			ErrorMsg: fmt.Sprintf("Failed to open result.json: %+v", err),
		})
		return
	}

	if processResult, ok := resultJson[pidStr]; !ok {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message":        "Process still processing",
			"videoGenerated": false,
		})
		return
	} else {
		if processResult.ProcessSucceed {
			c.JSON(http.StatusOK, map[string]interface{}{
				"outputPath":     outputPath,
				"videoGenerated": true,
				"exitCode":       processResult.ErrorCode,
			})
		} else {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"errorMsg":  fmt.Sprintf("Process PID %s failed to generate video", pidStr),
				"exitCode":  processResult.ErrorCode,
				"exitError": processResult.Error,
			})
		}
	}

}

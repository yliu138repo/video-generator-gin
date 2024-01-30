package videos

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h handler) DownloadVideo(c *gin.Context) {
	// Validate input
	filePathEncoded, ok := c.GetQuery("file_path")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: "file path not provided",
		})
		return
	}

	filePath, err := url.QueryUnescape(filePathEncoded)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{
			ErrorMsg: fmt.Sprintf("Failed to uncode the file_path parameter: %+v", err),
		})
		return
	}

	filePathParsed := filepath.Join(filePath)
	if _, err := os.Stat(filePathParsed); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			ErrorMsg: "file provided does not exist",
		})
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePathParsed))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePathParsed)
}

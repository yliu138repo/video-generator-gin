package books

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yliu138repo/video-generator-gin/pkg/common/models"
)

func (h handler) GetBooks(c *gin.Context) {
	var books []models.Book

	if result := h.DB.Find(&books); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)

		return
	}

	c.JSON(http.StatusOK, &books)
}

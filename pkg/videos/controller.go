package videos

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/videos")
	routes.POST("/", h.GenerateVideo)
	routes.POST("/cover", h.GenerateCoverPage)
	routes.GET("/", h.GetVideoStatus)
}

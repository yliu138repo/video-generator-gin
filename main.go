package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yliu138repo/video-generator-gin/docs"
	"github.com/yliu138repo/video-generator-gin/pkg/books"
	"github.com/yliu138repo/video-generator-gin/pkg/common/system"
	"github.com/yliu138repo/video-generator-gin/pkg/videos"
	"gorm.io/gorm"
)

//go:embed .env
var env string

func main() {
	if !system.CommandExists("ffmpeg") {
		log.Fatal("ffmpeg is not installed.")
	}

	viper.SetConfigType("env")
	viperReadErr := viper.ReadConfig(bytes.NewReader([]byte(env)))
	if viperReadErr != nil {
		log.Fatal("Failed to read env file.")
	}

	port := viper.Get("PORT").(string)
	// dbUrl := viper.Get("DB_URL").(string)

	r := gin.Default()
	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	// versioned API for backward compatibility
	v1 := r.Group(basePath)
	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// set middleware for gin
	m.Use(r)

	// h := db.Init(dbUrl)
	var h *gorm.DB

	books.RegisterRoutes(v1, h)
	videos.RegisterRoutes(v1, h)
	// register more routes here

	// Swagger routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(port)
}

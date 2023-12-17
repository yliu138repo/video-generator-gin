package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hellokvn/go-gin-api-medium/pkg/books"
	"github.com/hellokvn/go-gin-api-medium/pkg/common/db"
	"github.com/hellokvn/go-gin-api-medium/pkg/common/system"
	"github.com/hellokvn/go-gin-api-medium/pkg/videos"
	"github.com/spf13/viper"
)

func main() {
	if !system.CommandExists("ffmpeg") {
		log.Fatal("ffmpeg is not installed")
	}

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	port := viper.Get("PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	r := gin.Default()
	h := db.Init(dbUrl)

	books.RegisterRoutes(r, h)
	videos.RegisterRoutes(r, h)
	// register more routes here

	r.Run(port)
}

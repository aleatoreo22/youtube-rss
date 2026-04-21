package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"youtube-rss.api.aleatoreo.com/internal/api/handler"
	"youtube-rss.api.aleatoreo.com/internal/api/routes"
	"youtube-rss.api.aleatoreo.com/internal/repository"
	"youtube-rss.api.aleatoreo.com/internal/service"
)

var Handler *handler.Handler
var Service *service.Service

func main() {
	println("Hello World")

	dbType := "sqlite"
	repository, err := repository.New(dbType)
	if err != nil {
		panic(err)
	}
	Service = service.NewService(repository)
	Handler = handler.NewHandler(Service)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	routes.RegisterUserRoutes(r, Handler.UserHandler)
	routes.RegisterContentRoutes(r, Handler.ContentHandler)

	go getYoutubeRssContent()

	if err := r.Run(":1234"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func getYoutubeRssContent() {
	for {
		if !Handler.InSyncContent {
			Handler.InSyncContent = true
			Service.ContentService.SyncContent()
			Handler.InSyncContent = false
		}
		time.Sleep(time.Minute * 60)
	}
}

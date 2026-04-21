package routes

import (
	"github.com/gin-gonic/gin"
	"youtube-rss.api.aleatoreo.com/internal/api/handler"
)

func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	r.GET("/user/:id", h.GetById)
	r.PUT("/user/:id", h.Update)
	r.POST("/user", h.Create)

	r.POST("/user/:id/channel", h.AddUserYoutubeChannel)
	r.DELETE("/user/:id/channel/:channelId", h.DeleteUserChannel)
	r.GET("/user/:id/channel", h.GetUserChannels)

	r.GET("/user/:id/content/:date", h.GetContentByDate)
	r.GET("/user/:id/content", h.GetContentPaginated)
}

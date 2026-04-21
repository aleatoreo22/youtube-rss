package routes

import (
	"github.com/gin-gonic/gin"
	"youtube-rss.api.aleatoreo.com/internal/api/handler"
)

func RegisterContentRoutes(r *gin.Engine, h *handler.ContentHandler) {
	r.GET("/content/:date", h.GetByDate)
}

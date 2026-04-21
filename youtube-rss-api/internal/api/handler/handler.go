package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"youtube-rss.api.aleatoreo.com/internal/service"
)

type Handler struct {
	UserHandler    *UserHandler
	s              *service.Service
	InSyncContent  bool
	ContentHandler *ContentHandler
}

func NewHandler(s *service.Service) *Handler {
	h := &Handler{
		s: s,
	}
	user := UserHandler{h: h}
	h.UserHandler = &user

	content := ContentHandler{h: h}
	h.ContentHandler = &content
	return h
}

func (h *Handler) erroHandler(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "already exists") {
		code = http.StatusConflict
	}
	if strings.Contains(message, "not found") {
		code = http.StatusNotFound
	}

	println(err.Error())

	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}

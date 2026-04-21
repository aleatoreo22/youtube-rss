package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	h *Handler
}

func (ch *ContentHandler) GetByDate(c *gin.Context) {
	dateString := c.Param("date")
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		ch.h.erroHandler(c, err)
		return
	}

	contents, err := ch.h.s.ContentService.GetByDate(date)
	if err != nil {
		ch.h.erroHandler(c, err)
		return
	}

	if len(contents) == 0 {
		ch.h.erroHandler(c, errors.New("No content found"))
		return
	}

	c.JSON(http.StatusOK, contents)
}

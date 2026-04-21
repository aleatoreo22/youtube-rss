package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/internal/helper"
)

type UserHandler struct {
	h *Handler
}

func (u *UserHandler) GetContentByDate(c *gin.Context) {
	idString := c.Param("id")
	idNumber, err := strconv.Atoi(idString)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	dateString := c.Param("date")
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	content, err := u.h.s.UserService.GetContentByDate(idNumber, date)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, content)
}

func (u *UserHandler) GetContentPaginated(c *gin.Context) {
	idString := c.Param("id")
	idNumber, err := strconv.Atoi(idString)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	pagination, err := u.h.s.UserService.GetContentPaginated(idNumber, page, limit)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, pagination)
}

func (u *UserHandler) GetById(c *gin.Context) {
	idString := c.Param("id")
	idNumber, err := strconv.Atoi(idString)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	user, err := u.h.s.UserService.Get(idNumber)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	if user.Id == 0 {
		u.h.erroHandler(c, errors.New("User not found"))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Create(c *gin.Context) {
	user := domain.User{}
	c.BindJSON(&user)

	id, err := u.h.s.UserService.Create(user)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	user.Id = id
	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Update(c *gin.Context) {
	user, err := getUserAux(c, u)
	if err != nil || user.Id == 0 {
		return
	}
	id := user.Id

	err = c.BindJSON(&user)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	user.Id = id
	r, err := u.h.s.UserService.Upsert(user)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}
	if r == 0 {
		u.h.erroHandler(c, errors.New("User not found"))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) AddUserYoutubeChannel(c *gin.Context) {
	user, err := getUserAux(c, u)
	if err != nil || user.Id == 0 {
		return
	}

	channel := domain.Channel{}
	err = c.BindJSON(&channel)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	channel, err = helper.ValidateChannel(channel)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	id, err := u.h.s.ChannelService.Upsert(channel)
	channel.Id = id
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	userChannel, err := u.h.s.UserChannelService.GetByUserChannel(user.Id, channel.Id)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	if userChannel.Id > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "User already subscribed to this channel",
		})
	}

	_, err = u.h.s.UserChannelService.Upsert(domain.UserChannel{
		User:    user.Id,
		Channel: channel.Id,
	})
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (u *UserHandler) DeleteUserChannel(c *gin.Context) {
	idString := c.Param("id")
	idNumber, err := strconv.Atoi(idString)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	user, err := u.h.s.UserService.Get(idNumber)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	if user.Id == 0 {
		u.h.erroHandler(c, errors.New("User not found"))
		return
	}

	// Get all user channels
	userChannels, err := u.h.s.UserChannelService.GetByUser(user.Id)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	// Find and delete the channel
	for i, uc := range userChannels {
		if uc.Id == idNumber {
			_, err := u.h.s.UserChannelService.Delete(uc.Id)
			if err != nil {
				u.h.erroHandler(c, err)
				return
			}
			// Remove from list
			userChannels = append(userChannels[:i], userChannels[i+1:]...)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Channel removed successfully",
	})
}

func (u *UserHandler) GetUserChannels(c *gin.Context) {
	idString := c.Param("id")
	idNumber, err := strconv.Atoi(idString)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	user, err := u.h.s.UserService.Get(idNumber)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	if user.Id == 0 {
		u.h.erroHandler(c, errors.New("User not found"))
		return
	}

	// Get all user channels
	channels, err := u.h.s.UserChannelService.GetByUser(user.Id)
	if err != nil {
		u.h.erroHandler(c, err)
		return
	}

	// Convert to map for frontend
	// channels := make([]map[string]any, len(userChannels))
	// for i, uc := range userChannels {
	// 	channels[i] = map[string]any{
	// 		"id":      uc.Id,
	// 		"user":    uc.User,
	// 		"channel": uc.Channel,
	// 	}
	// }

	c.JSON(http.StatusOK, gin.H{
		"channels": channels,
	})
}

// Tlvz levar pro service dps
func getUserAux(c *gin.Context, u *UserHandler) (domain.User, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		u.h.erroHandler(c, err)
		return domain.User{}, err
	}

	user, err := u.h.s.UserService.Get(id)
	if err != nil {
		u.h.erroHandler(c, err)
		return domain.User{}, err
	}

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return domain.User{}, err
	}
	return user, nil
}

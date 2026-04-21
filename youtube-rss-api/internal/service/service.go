package service

import "youtube-rss.api.aleatoreo.com/internal/repository"

type Service struct {
	*UserService
	*ChannelService
	*UserChannelService
	*ContentService
	rp *repository.Repository
}

func NewService(rp *repository.Repository) *Service {
	s := &Service{
		rp: rp,
	}
	user := &UserService{s: s}
	channel := &ChannelService{s: s}
	userChannel := &UserChannelService{s: s}
	content := &ContentService{s: s}
	s.UserService = user
	s.ChannelService = channel
	s.UserChannelService = userChannel
	s.ContentService = content
	return s
}

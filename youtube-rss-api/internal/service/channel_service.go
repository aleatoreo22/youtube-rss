package service

import (
	"errors"

	"youtube-rss.api.aleatoreo.com/internal/domain"
)

type ChannelService struct {
	s *Service
}

func (cs *ChannelService) Get(id int) (domain.Channel, error) {
	return cs.s.rp.Channel.Get(id)
}

func (cs *ChannelService) Upsert(channel domain.Channel) (int, error) {
	return cs.s.rp.Channel.Upsert(channel)
}

func (cs *ChannelService) Create(channel domain.Channel) (int, error) {
	channel, err := cs.Get(channel.Id)
	if err != nil {
		return 0, err
	}
	if channel.Id > 0 {
		return 0, errors.New("Channel already exists")
	}

	return cs.Upsert(channel)
}

func (cs *ChannelService) GetAll() ([]domain.Channel, error) {
	return cs.s.rp.Channel.GetAll()
}

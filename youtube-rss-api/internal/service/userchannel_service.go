package service

import "youtube-rss.api.aleatoreo.com/internal/domain"

type UserChannelService struct {
	s *Service
}

func (ucs *UserChannelService) GetByUserChannel(user int, channel int) (domain.UserChannel, error) {
	return ucs.s.rp.UserChannel.GetByUserChannel(user, channel)
}

func (ucs *UserChannelService) GetByUser(user int) ([]domain.Channel, error) {
	return ucs.s.rp.UserChannel.GetByUser(user)
}

func (ucs *UserChannelService) Delete(id int) (int, error) {
	return ucs.s.rp.UserChannel.Delete(id)
}

func (ucs *UserChannelService) Upsert(userChannel domain.UserChannel) (int, error) {
	return ucs.s.rp.UserChannel.Upsert(userChannel)
}

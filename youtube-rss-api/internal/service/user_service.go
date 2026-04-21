package service

import (
	"errors"
	"time"

	"youtube-rss.api.aleatoreo.com/internal/domain"
)

type UserService struct {
	s *Service
}

func (us *UserService) GetContentByDate(id int, date time.Time) ([]domain.Content, error) {
	return us.s.rp.User.GetContentByDate(id, date)
}

func (us *UserService) GetContentPaginated(userId int, page, limit int) (domain.Pagination, error) {
	return us.s.rp.User.GetContentPaginated(userId, page, limit)
}

func (us *UserService) Get(id int) (domain.User, error) {
	return us.s.rp.User.Get(id)
}

func (us *UserService) Upsert(user domain.User) (int, error) {
	return us.s.rp.User.Upsert(user)
}

func (us *UserService) Create(user domain.User) (int, error) {
	userGet, err := us.Get(user.Id)
	if err != nil {
		return 0, err
	}
	if userGet.Id > 0 {
		return 0, errors.New("User already exists")
	}

	return us.Upsert(user)
}

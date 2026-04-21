package repository

import (
	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/package/database"
)

type Repository struct {
	database    *database.Database
	User        *UserRepository
	Channel     *ChannelRepository
	UserChannel *UserChannelRepository
	Content     *ContentRepository
}

func New(dbtype string) (*Repository, error) {
	database, err := database.New(dbtype)
	if err != nil {
		return nil, err
	}

	err = database.CreateDatabase([]any{
		domain.User{},
		domain.Channel{},
		domain.Content{},
		domain.UserChannel{},
		domain.UserContent{},
	})
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		database: database,
	}
	user := &UserRepository{repo: repo}
	channel := &ChannelRepository{repo: repo}
	userChannel := &UserChannelRepository{repo: repo}
	content := &ContentRepository{repo: repo}
	repo.User = user
	repo.Channel = channel
	repo.UserChannel = userChannel
	repo.Content = content
	return repo, nil
}

func (r *Repository) EasyExecute(query string) (int, error) {
	r.database.Framework.Open()
	result, err := r.database.Framework.Execute(query)
	r.database.Framework.Close()
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

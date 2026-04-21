package repository

import (
	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/package/database"
)

type ChannelRepository struct {
	repo *Repository
}

func (r *ChannelRepository) Upsert(channel domain.Channel) (int, error) {
	var query database.QueryBuilder

	query.WriteString("INSERT INTO Channel ( ")
	if channel.Id > 0 {
		query.WriteString(" id, ")
	}

	query.WriteString(" url, name, rssurl ) VALUES ( ")
	if channel.Id > 0 {
		query.WriteString(" ?id, ")
	}

	query.WriteString(" ?url, ?name, ?rssurl )  ")
	if channel.Id > 0 {
		query.WriteString("ON CONFLICT(id) DO UPDATE ")
	}

	query.AddParameter("?id", channel.Id)
	query.AddParameter("?url", channel.Url)
	query.AddParameter("?name", channel.Name)
	query.AddParameter("?rssurl", channel.RssUrl)

	return r.repo.EasyExecute(query.String())
}

func (r *ChannelRepository) Get(id int) (domain.Channel, error) {
	query := database.NewQueryBuilder("SELECT * FROM Channel WHERE id = ?id")

	query.AddParameter("?id", id)

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.Channel{})
	r.repo.database.Framework.Close()
	if err != nil {
		return domain.Channel{}, err
	}

	channels, err := database.ConvertQueryResult[domain.Channel](result)
	if err != nil {
		return domain.Channel{}, err
	}
	if len(channels) == 0 {
		return domain.Channel{}, nil
	}

	return channels[0], nil
}

func (r *ChannelRepository) GetAll() ([]domain.Channel, error) {
	query := database.NewQueryBuilder("SELECT * FROM Channel")

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.Channel{})
	r.repo.database.Framework.Close()
	if err != nil {
		return nil, err
	}

	channels, err := database.ConvertQueryResult[domain.Channel](result)
	if err != nil {
		return nil, err
	}
	if len(channels) == 0 {
		return nil, nil
	}

	return channels, nil
}

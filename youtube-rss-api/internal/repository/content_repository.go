package repository

import (
	"time"

	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/package/database"
)

type ContentRepository struct {
	repo *Repository
}

func (r *ContentRepository) GetByDate(date time.Time) ([]domain.Content, error) {
	query := database.NewQueryBuilder("SELECT * FROM Content WHERE Date > ?Date")
	query.AddParameter("?Date", date)

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.Content{})
	r.repo.database.Framework.Close()
	if err != nil {
		return nil, err
	}

	contents, err := database.ConvertQueryResult[domain.Content](result)

	if err != nil {
		return nil, err
	}
	if len(contents) == 0 {
		return nil, nil
	}

	return contents, nil
}

func (r *ContentRepository) GetByUrl(url string) (domain.Content, error) {
	query := database.NewQueryBuilder("SELECT * FROM Content WHERE Url = ?Url")
	query.AddParameter("?Url", url)

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.Content{})
	r.repo.database.Framework.Close()
	if err != nil {
		return domain.Content{}, err
	}

	contents, err := database.ConvertQueryResult[domain.Content](result)
	if err != nil {
		return domain.Content{}, err
	}
	if len(contents) == 0 {
		return domain.Content{}, nil
	}

	return contents[0], nil
}

func (r *ContentRepository) Upsert(content domain.Content) (int, error) {
	var query database.QueryBuilder

	query.WriteString("INSERT INTO Content ( ")
	query.WriteString(" Channel, Url, Title, Image, Date ) VALUES ( ")
	query.WriteString(" ?Channel, ?Url, ?Title, ?Image, ?Date ) ")
	query.WriteString("ON CONFLICT(Url) DO UPDATE SET Title = excluded.Title, ")
	query.WriteString(" Image = excluded.Image, Date = excluded.Date ")
	query.AddParameter("?Channel", content.Channel)
	query.AddParameter("?Url", content.Url)
	query.AddParameter("?Title", content.Title)
	query.AddParameter("?Image", content.Image)
	query.AddParameter("?Date", content.Date)

	return r.repo.EasyExecute(query.String())
}

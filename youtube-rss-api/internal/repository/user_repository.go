package repository

import (
	"time"

	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/package/database"
)

type UserRepository struct {
	repo *Repository
}

func (r *UserRepository) GetContentByDate(id int, date time.Time) ([]domain.Content, error) {
	query := database.NewQueryBuilder("SELECT Content.* FROM User ")
	query.WriteString(" INNER JOIN UserChannel ON User.Id = UserChannel.User ")
	query.WriteString(" INNER JOIN Content ON Content.Channel = UserChannel.Channel ")
	query.WriteString(" WHERE User.id = ?id AND Content.date > ?date")
	query.WriteString(" ORDER BY Content.Date DESC")

	query.AddParameter("?id", id)
	query.AddParameter("?date", date)

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

	return contents, nil
}

func (r *UserRepository) GetContentPaginated(userId int, page, limit int) (domain.Pagination, error) {
	offset := (page - 1) * limit

	countQuery := database.NewQueryBuilder("SELECT COUNT(*) Count FROM User ")
	countQuery.WriteString(" INNER JOIN UserChannel ON User.Id = UserChannel.User ")
	countQuery.WriteString(" INNER JOIN Content ON Content.Channel = UserChannel.Channel ")
	countQuery.WriteString(" WHERE User.id = ?id")

	countQuery.AddParameter("?id", userId)

	r.repo.database.Framework.Open()
	countResult, err := r.repo.database.Framework.Query(countQuery.String(), struct{ Count int }{})
	r.repo.database.Framework.Close()
	if err != nil {
		return domain.Pagination{}, err
	}

	var total int
	if len(countResult) > 0 {
		total = int(countResult[0].(struct{ Count int }).Count)
	}

	query := database.NewQueryBuilder("SELECT Content.* FROM User ")
	query.WriteString(" INNER JOIN UserChannel ON User.Id = UserChannel.User ")
	query.WriteString(" INNER JOIN Content ON Content.Channel = UserChannel.Channel ")
	query.WriteString(" WHERE User.id = ?id")
	query.WriteString(" ORDER BY Content.Date DESC")
	query.WriteString(" LIMIT ?limit OFFSET ?offset")

	query.AddParameter("?id", userId)
	query.AddParameter("?limit", limit)
	query.AddParameter("?offset", offset)

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.Content{})
	r.repo.database.Framework.Close()
	if err != nil {
		return domain.Pagination{}, err
	}

	contents, err := database.ConvertQueryResult[domain.Content](result)
	if err != nil {
		return domain.Pagination{}, err
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return domain.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Items:      contents,
	}, nil
}

func (r *UserRepository) Get(id int) (domain.User, error) {
	query := database.NewQueryBuilder("SELECT * FROM User WHERE id = ?id")
	query.AddParameter("?id", id)

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.User{})
	r.repo.database.Framework.Close()
	if err != nil {
		return domain.User{}, err
	}

	users, err := database.ConvertQueryResult[domain.User](result)
	if err != nil {
		return domain.User{}, err
	}
	if len(users) == 0 {
		return domain.User{}, nil
	}

	return users[0], nil
}

func (r *UserRepository) Upsert(user domain.User) (int, error) {
	var query database.QueryBuilder

	query.WriteString("INSERT INTO User ( ")
	if user.Id > 0 {
		query.WriteString(" id, ")
	}
	query.WriteString(" username ) VALUES ( ")
	if user.Id > 0 {
		query.WriteString(" ?id, ")
	}
	query.WriteString(" ?username ) ")
	if user.Id > 0 {
		query.WriteString("ON CONFLICT(id) DO UPDATE SET username = excluded.username")
	}

	println("Parameter Username: " + user.Username)
	query.AddParameter("?username", user.Username)
	query.AddParameter("?id", user.Id)

	println(query.String())

	return r.repo.EasyExecute(query.String())
}

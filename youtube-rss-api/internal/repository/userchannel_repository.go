package repository

import (
	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/package/database"
)

type UserChannelRepository struct {
	repo *Repository
}

func (r *UserChannelRepository) GetByUserChannel(userId int, channelId int) (domain.UserChannel, error) {
	query := database.NewQueryBuilder("SELECT * FROM UserChannel WHERE user = ?userId AND channel = ?channelId")

	query.AddParameter("?userId", userId)
	query.AddParameter("?channelId", channelId)

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.UserChannel{})
	r.repo.database.Framework.Close()
	if err != nil {
		return domain.UserChannel{}, err
	}

	userChannels, err := database.ConvertQueryResult[domain.UserChannel](result)
	if err != nil {
		return domain.UserChannel{}, err
	}
	if len(userChannels) == 0 {
		return domain.UserChannel{}, nil
	}

	return userChannels[0], nil
}

func (r *UserChannelRepository) GetByUser(userId int) ([]domain.Channel, error) {
	query := database.NewQueryBuilder("SELECT c.* FROM UserChannel uc ")
	query.WriteString(" INNER JOIN Channel c on uc.channel = c.id ")

	query.WriteString(" WHERE uc.user = ?userId")

	query.AddParameter("?userId", userId)

	r.repo.database.Framework.Open()
	result, err := r.repo.database.Framework.Query(query.String(), domain.Channel{})
	r.repo.database.Framework.Close()
	if err != nil {
		return nil, err
	}

	userChannels, err := database.ConvertQueryResult[domain.Channel](result)
	if err != nil {
		return nil, err
	}
	if len(userChannels) == 0 {
		return nil, nil
	}

	return userChannels, nil
}

func (r *UserChannelRepository) Upsert(userChannel domain.UserChannel) (int, error) {
	var query database.QueryBuilder

	query.WriteString("INSERT INTO UserChannel ( ")
	if userChannel.Id > 0 {
		query.WriteString(" id, ")
	}

	query.WriteString(" user, channel ) VALUES ( ")
	if userChannel.Id > 0 {
		query.WriteString(" ?id, ")
	}

	query.WriteString(" ?user, ?channel )  ")
	if userChannel.Id > 0 {
		query.WriteString("ON CONFLICT(id) DO UPDATE ")
	}

	query.AddParameter("?id", userChannel.Id)
	query.AddParameter("?user", userChannel.User)
	query.AddParameter("?channel", userChannel.Channel)

	return r.repo.EasyExecute(query.String())
}

func (r *UserChannelRepository) Delete(id int) (int, error) {
	query := database.NewQueryBuilder("DELETE FROM UserChannel WHERE id = ?id")

	query.AddParameter("?id", id)

	return r.repo.EasyExecute(query.String())
}

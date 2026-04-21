package service

import (
	"time"

	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/package/youtuberss"
)

type ContentService struct {
	s *Service
}

func (cs *ContentService) GetByDate(date time.Time) ([]domain.Content, error) {
	return cs.s.rp.Content.GetByDate(date)
}

func (cs *ContentService) SyncContent() error {
	channels, err := cs.s.ChannelService.GetAll()
	if err != nil {
		return err
	}

	ytrss := youtuberss.New()
	var contents []domain.Content
	for _, channel := range channels {
		feed, err := ytrss.GetFeed(channel.RssUrl)
		if err != nil {
			return err
		}
		for _, entry := range feed.Entry {
			date, err := time.Parse(time.RFC3339, entry.Published)
			if err != nil {
				return err
			}

			contents = append(contents, domain.Content{
				Channel: channel.Id,
				Url:     entry.Group.Content.URL,
				Title:   entry.Title,
				Image:   entry.Group.Thumbnail.URL,
				Date:    date,
			})
		}
	}

	for _, content := range contents {
		_, err := cs.s.rp.Content.Upsert(content)
		if err != nil {
			return err
		}
	}
	return nil
}

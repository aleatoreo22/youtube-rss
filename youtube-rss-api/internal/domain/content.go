package domain

import "time"

type Content struct {
	Id      int       `sql:"key" json:"id"`
	Url     string    `sql:"strlen 255, unique" json:"url"`
	Channel int       `json:"channel"`
	Title   string    `sql:"strlen 255" json:"title"`
	Image   string    `sql:"strlen 255" json:"image"`
	Date    time.Time `json:"date"`
}

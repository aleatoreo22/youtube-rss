package domain

type Channel struct {
	Id     int    `sql:"key" json:"id"`
	Url    string `sql:"strlen 255" json:"url"`
	Name   string `sql:"strlen 255" json:"name"`
	RssUrl string `sql:"strlen 255" json:"rss_url"`
}

package domain

type UserChannel struct {
	Id      int `sql:"key" json:"id"`
	User    int `json:"user"`
	Channel int `json:"channel"`
}

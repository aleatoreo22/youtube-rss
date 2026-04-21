package domain

type User struct {
	Id       int    `sql:"key" json:"id"`
	Username string `sql:"strlen 255" json:"username"`
}

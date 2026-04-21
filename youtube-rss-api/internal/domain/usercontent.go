package domain

type UserContent struct {
	Id      int  `sql:"key" json:"id"`
	User    int  `json:"user"`
	Content int  `json:"channel"`
	Watched bool `json:"watched"`
	Hidden  bool `json:"hidden"`
	Time    int  `json:"time"`
}

type Pagination struct {
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	Total      int       `json:"total"`
	TotalPages int       `json:"totalPages"`
	Items      []Content `json:"items"`
}

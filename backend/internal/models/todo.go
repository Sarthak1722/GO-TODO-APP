package models

type Todo struct {
	ID        int    `json:"id"`
	UserID    string `json:"-"`
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}

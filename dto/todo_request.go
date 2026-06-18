package dto

type CreateTodoRequest struct {
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}
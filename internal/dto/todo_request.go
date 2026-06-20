package dto

type CreateTodoRequest struct {
	Body      string `json:"body" validate:"required,min=3,max=100"`
	Completed bool   `json:"completed"`
}

type PatchTodoRequest struct {
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}

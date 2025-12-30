package model

import "github.com/Tortik3000/todo-list/internal/model"

type TaskRequest struct {
	Header      string `json:"header" example:"Buy milk"`
	Description string `json:"description" example:"Go to the store and buy milk"`
	Completed   bool   `json:"completed" example:"false"`
}

type TaskResponse struct {
	ID          int64  `json:"id" example:"1"`
	Header      string `json:"header" example:"Buy milk"`
	Description string `json:"description" example:"Go to the store and buy milk"`
	Completed   bool   `json:"completed" example:"false"`
}

func ModelToResponse(t model.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Header:      t.Header,
		Description: t.Description,
		Completed:   t.Completed,
	}
}

func RequestToModel(t TaskRequest) model.Task {
	return model.Task{
		Header:      t.Header,
		Description: t.Description,
		Completed:   t.Completed,
	}
}

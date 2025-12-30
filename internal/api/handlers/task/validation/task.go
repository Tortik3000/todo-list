package validation

import (
	"github.com/Tortik3000/todo-list/internal/api/handlers/task/model"
)

func ValidateTaskRequest(req model.TaskRequest) error {
	if req.Header == "" {
		return model.ErrEmptyHeader
	}
	return nil
}

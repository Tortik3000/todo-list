package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	handlerModel "github.com/Tortik3000/todo-list/internal/api/handlers/task/model"
	domainModel "github.com/Tortik3000/todo-list/internal/model"
)

func TestModelToResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input domainModel.Task
	}{
		{
			name: "success convert",
			input: domainModel.Task{
				ID:          100,
				Header:      "header",
				Description: "description",
				Completed:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := handlerModel.ModelToResponse(tt.input)

			assert.Equal(t, tt.input.ID, result.ID)
			assert.Equal(t, tt.input.Header, result.Header)
			assert.Equal(t, tt.input.Description, result.Description)
			assert.Equal(t, tt.input.Completed, result.Completed)
		})
	}
}

func TestRequestToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input handlerModel.TaskRequest
	}{
		{
			name: "success convert",
			input: handlerModel.TaskRequest{
				Header:      "header",
				Description: "description",
				Completed:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := handlerModel.RequestToModel(tt.input)

			assert.Equal(t, tt.input.Header, result.Header)
			assert.Equal(t, tt.input.Description, result.Description)
			assert.Equal(t, tt.input.Completed, result.Completed)
			assert.Equal(t, int64(0), result.ID)
		})
	}
}

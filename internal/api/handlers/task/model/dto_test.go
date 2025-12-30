package model_test

import (
	"github.com/Tortik3000/todo-list/pkg/tests"
	"testing"

	handlerModel "github.com/Tortik3000/todo-list/internal/api/handlers/task/model"
	domainModel "github.com/Tortik3000/todo-list/internal/model"
)

func TestModelToResponse(t *testing.T) {
	t.Parallel()

	testsCases := []struct {
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

	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := handlerModel.ModelToResponse(tt.input)

			tests.AssertEqual(t, tt.input.ID, result.ID)
			tests.AssertEqual(t, tt.input.Header, result.Header)
			tests.AssertEqual(t, tt.input.Description, result.Description)
			tests.AssertEqual(t, tt.input.Completed, result.Completed)
		})
	}
}

func TestRequestToModel(t *testing.T) {
	t.Parallel()

	testsCases := []struct {
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

	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := handlerModel.RequestToModel(tt.input)

			tests.AssertEqual(t, tt.input.Header, result.Header)
			tests.AssertEqual(t, tt.input.Description, result.Description)
			tests.AssertEqual(t, tt.input.Completed, result.Completed)
			tests.AssertEqual(t, int64(0), result.ID)
		})
	}
}

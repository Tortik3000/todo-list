package validation_test

import (
	"testing"

	"github.com/Tortik3000/todo-list/pkg/tests"

	"github.com/Tortik3000/todo-list/internal/api/handlers/task/model"
	"github.com/Tortik3000/todo-list/internal/api/handlers/task/validation"
)

func Test_ValidateTaskRequest(t *testing.T) {
	t.Parallel()

	testsCases := []struct {
		name string
		req  model.TaskRequest
		err  error
	}{
		{
			name: "success valid task",
			req: model.TaskRequest{
				Header:      "header",
				Description: "description",
			},
			err: nil,
		},
		{
			name: "fail empty header",
			req: model.TaskRequest{
				Header:      "",
				Description: "description",
			},
			err: model.ErrEmptyHeader,
		},
	}

	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validation.ValidateTaskRequest(tt.req)
			if tt.err != nil {
				tests.RequireErrorIs(t, err, tt.err)
			} else {
				tests.RequireNoError(t, err)
			}
		})
	}
}

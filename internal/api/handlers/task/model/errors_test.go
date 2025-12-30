package model_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	handlerModel "github.com/Tortik3000/todo-list/internal/api/handlers/task/model"
	modelErr "github.com/Tortik3000/todo-list/internal/model/error"
)

func TestToHttpErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputErr       error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "map ErrTaskNotFound to 404",
			inputErr:       modelErr.ErrTaskNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "map ErrEmptyHeader to 400",
			inputErr:       handlerModel.ErrEmptyHeader,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "map ErrInvalidJSON to 400",
			inputErr:       handlerModel.ErrInvalidJSON,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "map ErrInvalidID to 400",
			inputErr:       handlerModel.ErrInvalidID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "map wrapped ErrTaskNotFound to 404",
			inputErr:       fmt.Errorf("wrapped: %w", modelErr.ErrTaskNotFound),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "map wrapped ErrEmptyHeader to 400",
			inputErr:       fmt.Errorf("validation failed: %w", handlerModel.ErrEmptyHeader),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "map unknown error to 500",
			inputErr:       errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			handlerModel.ToHttpErr(w, tt.inputErr)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

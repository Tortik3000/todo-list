package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpUtils "github.com/Tortik3000/todo-list/pkg/http"
)

func TestParseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pathValue string
		wantID    int64
		wantErr   bool
	}{
		{
			name:      "success positive id",
			pathValue: "123",
			wantID:    123,
			wantErr:   false,
		},
		{
			name:      "fail empty id",
			pathValue: "",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "fail non-numeric id",
			pathValue: "abc",
			wantID:    0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
			req.SetPathValue("id", tt.pathValue)

			gotID, err := httpUtils.ParseID(req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantID, gotID)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name       string
		status     int
		payload    any
		wantBody   string
		wantStatus int
	}{
		{
			name:       "success struct",
			status:     http.StatusOK,
			payload:    testStruct{Name: "Alice", Age: 30},
			wantBody:   `{"name":"Alice","age":30}` + "\n",
			wantStatus: http.StatusOK,
		},
		{
			name:       "success map",
			status:     http.StatusCreated,
			payload:    map[string]string{"foo": "bar"},
			wantBody:   `{"foo":"bar"}` + "\n",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "success nil payload",
			status:     http.StatusNoContent,
			payload:    nil,
			wantBody:   "null\n",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			httpUtils.WriteJSON(w, tt.status, tt.payload)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

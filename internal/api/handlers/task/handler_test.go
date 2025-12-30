package task_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	taskHandler "github.com/Tortik3000/todo-list/internal/api/handlers/task"
	handlerModel "github.com/Tortik3000/todo-list/internal/api/handlers/task/model"
	taskRepo "github.com/Tortik3000/todo-list/internal/repository/in_memory/task"
	taskService "github.com/Tortik3000/todo-list/internal/service/task"
)

func TestHandler_CreateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		task0 handlerModel.TaskRequest
		task1 handlerModel.TaskRequest
	}{
		{
			name: "success create task",
			task0: handlerModel.TaskRequest{
				Header:      "header_0",
				Description: "description_0",
				Completed:   false,
			},
			task1: handlerModel.TaskRequest{
				Header:    "header_1",
				Completed: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newTestHandler()

			created0 := createTaskDirect(t, h, tt.task0)
			assert.Equal(t, int64(0), created0.ID)
			assert.Equal(t, tt.task0.Header, created0.Header)
			assert.Equal(t, tt.task0.Description, created0.Description)
			assert.Equal(t, tt.task0.Completed, created0.Completed)

			created1 := createTaskDirect(t, h, tt.task1)
			assert.Equal(t, int64(1), created1.ID)
			assert.Equal(t, tt.task1.Header, created1.Header)
			assert.Equal(t, tt.task1.Description, created1.Description)
			assert.Equal(t, tt.task1.Completed, created1.Completed)
		})
	}
}

func TestHandler_GetTaskByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		taskID string
		status int
	}{
		{
			name:   "success get task",
			taskID: "0",
			status: http.StatusOK,
		},
		{
			name:   "fail get not found",
			taskID: "1",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newTestHandler()

			expected := createTaskDirect(t, h, handlerModel.TaskRequest{
				Header:      "header",
				Description: "description",
				Completed:   false,
			})

			req := httptest.NewRequest(http.MethodGet, "/todos/"+tt.taskID, nil)
			req = withIDParam(req, tt.taskID)
			rr := httptest.NewRecorder()

			h.GetTaskByID(rr, req)

			require.Equal(t, tt.status, rr.Code)

			if tt.status == http.StatusOK {
				got := handlerModel.TaskResponse{}
				require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
				assert.Equal(t, expected, got)
			}
		})
	}
}

func TestHandler_GetAllTasks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup []handlerModel.TaskRequest
	}{
		{
			name: "success get tasks",
			setup: []handlerModel.TaskRequest{
				{Header: "h0", Description: "d0", Completed: false},
				{Header: "h1", Description: "d1", Completed: true},
			},
		},
		{
			name:  "success get empty tasks",
			setup: []handlerModel.TaskRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newTestHandler()

			expected := make([]handlerModel.TaskResponse, 0, len(tt.setup))
			for _, req := range tt.setup {
				expected = append(expected, createTaskDirect(t, h, req))
			}

			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rr := httptest.NewRecorder()

			h.GetAllTasks(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			got := make([]handlerModel.TaskResponse, 0)
			require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))

			sort.Slice(got, func(i, j int) bool { return got[i].ID < got[j].ID })
			sort.Slice(expected, func(i, j int) bool { return expected[i].ID < expected[j].ID })

			assert.Equal(t, expected, got)
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		taskID string
		status int
	}{
		{
			name:   "success update task",
			taskID: "0",
			status: http.StatusOK,
		},
		{
			name:   "not found task for update",
			taskID: "1",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newTestHandler()
			_ = createTaskDirect(t, h, handlerModel.TaskRequest{
				Header:      "old",
				Description: "old_desc",
				Completed:   false,
			})

			updateReq := handlerModel.TaskRequest{
				Header:      "new",
				Description: "new_desc",
				Completed:   true,
			}
			body, err := json.Marshal(updateReq)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/todos/"+tt.taskID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req = withIDParam(req, tt.taskID)
			rr := httptest.NewRecorder()

			h.UpdateTask(rr, req)
			require.Equal(t, tt.status, rr.Code)

			if tt.status != http.StatusOK {
				return
			}

			getReq := httptest.NewRequest(http.MethodGet, "/todos/"+tt.taskID, nil)
			getReq = withIDParam(getReq, tt.taskID)
			getRR := httptest.NewRecorder()

			h.GetTaskByID(getRR, getReq)
			require.Equal(t, http.StatusOK, getRR.Code)

			got := handlerModel.TaskResponse{}
			require.NoError(t, json.NewDecoder(getRR.Body).Decode(&got))

			assert.Equal(t, int64(0), got.ID)
			assert.Equal(t, updateReq.Header, got.Header)
			assert.Equal(t, updateReq.Description, got.Description)
			assert.Equal(t, updateReq.Completed, got.Completed)
		})
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		taskID string
		status int
	}{
		{
			name:   "success delete task",
			taskID: "0",
			status: http.StatusOK,
		},
		{
			name:   "not found task for delete",
			taskID: "1",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newTestHandler()
			_ = createTaskDirect(t, h, handlerModel.TaskRequest{
				Header:      "h",
				Description: "d",
				Completed:   false,
			})

			req := httptest.NewRequest(http.MethodDelete, "/todos/"+tt.taskID, nil)
			req = withIDParam(req, tt.taskID)
			rr := httptest.NewRecorder()

			h.DeleteTask(rr, req)
			require.Equal(t, tt.status, rr.Code)

			if tt.status != http.StatusOK {
				return
			}
			getReq := httptest.NewRequest(http.MethodGet, "/todos/"+tt.taskID, nil)
			getReq = withIDParam(getReq, tt.taskID)
			getRR := httptest.NewRecorder()

			h.GetTaskByID(getRR, getReq)
			require.Equal(t, http.StatusNotFound, getRR.Code)
		})
	}
}

func newTestHandler() taskHandler.Handler {
	repository := taskRepo.New()
	service := taskService.New(repository)
	return taskHandler.New(service)
}

func createTaskDirect(t *testing.T, h taskHandler.Handler, reqBody handlerModel.TaskRequest) handlerModel.TaskResponse {
	t.Helper()

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.CreateTask(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	resp := handlerModel.TaskResponse{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))

	return resp
}

func withIDParam(r *http.Request, id string) *http.Request {
	r.SetPathValue("id", id)
	return r
}

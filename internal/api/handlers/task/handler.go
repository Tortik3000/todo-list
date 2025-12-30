package task

import (
	"context"
	"encoding/json"
	"net/http"

	handlerModel "github.com/Tortik3000/todo-list/internal/api/handlers/task/model"
	"github.com/Tortik3000/todo-list/internal/api/handlers/task/validation"
	domainModel "github.com/Tortik3000/todo-list/internal/model"
	httpUtils "github.com/Tortik3000/todo-list/pkg/http"
)

type Handler interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetAllTasks(w http.ResponseWriter, r *http.Request)
	GetTaskByID(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

type (
	taskService interface {
		CreateTask(ctx context.Context, task domainModel.Task) (*domainModel.Task, error)
		GetAllTasks(ctx context.Context) ([]domainModel.Task, error)
		GetTaskByID(ctx context.Context, id int64) (*domainModel.Task, error)
		UpdateTask(ctx context.Context, task domainModel.Task) error
		DeleteTask(ctx context.Context, id int64) error
	}
)

type handler struct {
	taskService taskService
}

func New(
	taskService taskService,
) *handler {
	return &handler{
		taskService: taskService,
	}
}

var _ Handler = (*handler)(nil)

// CreateTask godoc
// @Summary Create a new task
// @Param task body handlerModel.TaskRequest true "Task to create"
// @Success 201 {object} handlerModel.TaskResponse
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal server error"
// @Router /todos [post]
func (h handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			handlerModel.ToHttpErr(w, err)
		}
	}()

	req := handlerModel.TaskRequest{}
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = handlerModel.ErrInvalidJSON
		return
	}

	if err = validation.ValidateTaskRequest(req); err != nil {
		return
	}

	task := handlerModel.RequestToModel(req)
	createdTask, err := h.taskService.CreateTask(r.Context(), task)
	if err != nil {
		return
	}

	response := handlerModel.ModelToResponse(*createdTask)
	httpUtils.WriteJSON(w, http.StatusCreated, response)
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Success 200 {array} handlerModel.TaskResponse
// @Failure 500 {object} string "Internal server error"
// @Router /todos [get]
func (h handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			handlerModel.ToHttpErr(w, err)
		}
	}()

	tasks, err := h.taskService.GetAllTasks(r.Context())
	if err != nil {
		return
	}

	response := make([]handlerModel.TaskResponse, len(tasks))
	for i, t := range tasks {
		response[i] = handlerModel.ModelToResponse(t)
	}

	httpUtils.WriteJSON(w, http.StatusOK, response)
}

// GetTaskByID godoc
// @Summary Get a task by ID
// @Param id path int true "Task ID"
// @Success 200 {object} handlerModel.TaskResponse
// @Failure 400 {object} string "Invalid ID"
// @Failure 404 {object} string "Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /todos/{id} [get]
func (h handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			handlerModel.ToHttpErr(w, err)
		}
	}()

	id, err := httpUtils.ParseID(r)
	if err != nil {
		err = handlerModel.ErrInvalidID
		return
	}

	task, err := h.taskService.GetTaskByID(r.Context(), id)
	if err != nil {
		return
	}

	response := handlerModel.ModelToResponse(*task)
	httpUtils.WriteJSON(w, http.StatusOK, response)
}

// UpdateTask godoc
// @Summary Update an existing task
// @Param id path int true "Task ID"
// @Param task body handlerModel.TaskRequest true "Updated task details"
// @Success 200 {object} string "Success"
// @Failure 400 {object} string "Invalid input or ID"
// @Failure 404 {object} string "Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /todos/{id} [put]
func (h handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			handlerModel.ToHttpErr(w, err)
		}
	}()

	id, err := httpUtils.ParseID(r)
	if err != nil {
		err = handlerModel.ErrInvalidID
		return
	}

	var req handlerModel.TaskRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = handlerModel.ErrInvalidJSON
		return
	}

	if err = validation.ValidateTaskRequest(req); err != nil {
		return
	}

	task := handlerModel.RequestToModel(req)
	task.ID = id

	err = h.taskService.UpdateTask(r.Context(), task)
	if err != nil {
		return
	}

	httpUtils.WriteJSON(w, http.StatusOK, "")
}

// DeleteTask godoc
// @Summary Delete a task
// @Param id path int true "Task ID"
// @Success 200 {object} string "Success"
// @Failure 400 {object} string "Invalid ID"
// @Failure 404 {object} string "Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /todos/{id} [delete]
func (h handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			handlerModel.ToHttpErr(w, err)
		}
	}()

	id, err := httpUtils.ParseID(r)
	if err != nil {
		err = handlerModel.ErrInvalidID
		return
	}

	err = h.taskService.DeleteTask(r.Context(), id)
	if err != nil {
		return
	}

	httpUtils.WriteJSON(w, http.StatusOK, "")
}

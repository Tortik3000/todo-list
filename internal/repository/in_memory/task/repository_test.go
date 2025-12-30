package task_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Tortik3000/todo-list/internal/model"
	modelErr "github.com/Tortik3000/todo-list/internal/model/error"
	repo "github.com/Tortik3000/todo-list/internal/repository/in_memory/task"
)

func TestRepository_CreateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		task0 model.Task
		task1 model.Task
	}{
		{
			name: "success create task",
			task0: model.Task{
				ID:          0,
				Header:      "header_0",
				Description: "description_0",
				Completed:   false,
			},
			task1: model.Task{
				ID:        1,
				Header:    "header_1",
				Completed: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := repo.New()

			res, err := repository.CreateTask(t.Context(), tt.task0)
			require.NoError(t, err)
			taskEqual(t, tt.task0, *res)

			res, err = repository.CreateTask(t.Context(), tt.task1)
			require.NoError(t, err)
			taskEqual(t, tt.task1, *res)
		})
	}
}

func TestRepository_GetTaskByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		taskID int64
		err    error
	}{
		{
			name:   "success get task",
			taskID: 0,
			err:    nil,
		},
		{
			name:   "fail get non-existing task",
			taskID: 1,
			err:    modelErr.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := repo.New()
			task := setupTask(t, repository, generateTask())

			id := tt.taskID
			res, err := repository.GetTaskByID(t.Context(), id)
			if tt.err != nil {
				require.ErrorIs(t, tt.err, err)
			} else {
				require.NoError(t, err)
				taskEqual(t, task, *res)
			}
		})
	}
}

func TestRepository_GetAllTasks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tasks []model.Task
	}{
		{
			name:  "success get tasks",
			tasks: []model.Task{generateTask(), generateTask()},
		},
		{
			name:  "success get empty tasks",
			tasks: []model.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := repo.New()
			for i, task := range tt.tasks {
				tt.tasks[i] = setupTask(t, repository, task)
			}

			res, err := repository.GetAllTasks(t.Context())
			require.NoError(t, err)
			sort.Slice(res, func(i, j int) bool {
				return res[i].ID < res[j].ID
			})

			assert.Equal(t, tt.tasks, res)
		})
	}
}

func TestRepository_UpdateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		taskID int64
		err    error
	}{
		{
			name:   "success update task",
			taskID: 0,
			err:    nil,
		},
		{
			name:   "not found task for update",
			taskID: 1,
			err:    modelErr.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := repo.New()
			setupTask(t, repository, generateTask())
			newTask := generateTask()
			newTask.ID = tt.taskID

			err := repository.UpdateTask(t.Context(), newTask)
			if tt.err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)

			getTask, err := repository.GetTaskByID(t.Context(), newTask.ID)
			require.NoError(t, err)
			taskEqual(t, newTask, *getTask)
		})
	}
}

func TestRepository_DeleteTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		taskID int64
		err    error
	}{
		{
			name:   "success delete task",
			taskID: 0,
			err:    nil,
		},
		{
			name:   "not found task for delete",
			taskID: 1,
			err:    modelErr.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := repo.New()
			setupTask(t, repository, generateTask())

			err := repository.DeleteTask(t.Context(), tt.taskID)
			if tt.err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)

			_, err = repository.GetTaskByID(t.Context(), tt.taskID)
			require.ErrorIs(t, err, modelErr.ErrTaskNotFound)
		})
	}
}

func taskEqual(t *testing.T, expected, actual model.Task) {
	t.Helper()

	assert.Equal(t, expected.Header, actual.Header)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.Completed, actual.Completed)
}

func setupTask(t *testing.T, repository repo.Repository, inTask model.Task) model.Task {
	t.Helper()

	task, err := repository.CreateTask(t.Context(), inTask)
	require.NoError(t, err)

	return *task
}

func generateTask() model.Task {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return model.Task{
		Header:      fmt.Sprintf("header %d", r.Intn(10)),
		Description: fmt.Sprintf("description %d", r.Intn(10)),
		Completed:   r.Intn(2) == 1,
	}
}

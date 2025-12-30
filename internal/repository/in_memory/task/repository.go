package task

import (
	"context"
	"sync"

	domainModel "github.com/Tortik3000/todo-list/internal/model"
	modelErr "github.com/Tortik3000/todo-list/internal/model/error"
)

type Repository interface {
	CreateTask(ctx context.Context, task domainModel.Task) (*domainModel.Task, error)
	GetAllTasks(ctx context.Context) ([]domainModel.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*domainModel.Task, error)
	UpdateTask(ctx context.Context, task domainModel.Task) error
	DeleteTask(ctx context.Context, id int64) error
}

type repository struct {
	tasks  map[int64]*domainModel.Task
	mu     *sync.RWMutex
	nextID int64
}

var _ Repository = (*repository)(nil)

func New() *repository {
	return &repository{
		tasks:  make(map[int64]*domainModel.Task),
		mu:     &sync.RWMutex{},
		nextID: 0,
	}
}

func (r *repository) CreateTask(
	_ context.Context,
	task domainModel.Task,
) (*domainModel.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.nextID++

	newTask := &domainModel.Task{
		ID:          task.ID,
		Header:      task.Header,
		Description: task.Description,
		Completed:   task.Completed,
	}
	r.tasks[task.ID] = newTask

	return newTask, nil
}

func (r *repository) GetAllTasks(
	_ context.Context,
) ([]domainModel.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]domainModel.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		tasks = append(tasks, *t)
	}

	return tasks, nil
}

func (r *repository) GetTaskByID(
	_ context.Context,
	id int64,
) (*domainModel.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, modelErr.ErrTaskNotFound
	}

	return task, nil
}

func (r *repository) UpdateTask(
	_ context.Context,
	task domainModel.Task,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[task.ID]; !ok {
		return modelErr.ErrTaskNotFound
	}

	r.tasks[task.ID] = &task
	return nil
}

func (r *repository) DeleteTask(
	_ context.Context,
	id int64,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return modelErr.ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

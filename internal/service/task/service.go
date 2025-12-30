package task

import (
	"context"

	domainModel "github.com/Tortik3000/todo-list/internal/model"
)

type Service interface {
	CreateTask(ctx context.Context, task domainModel.Task) (*domainModel.Task, error)
	GetAllTasks(ctx context.Context) ([]domainModel.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*domainModel.Task, error)
	UpdateTask(ctx context.Context, task domainModel.Task) error
	DeleteTask(ctx context.Context, id int64) error
}

type (
	taskRepo interface {
		CreateTask(ctx context.Context, task domainModel.Task) (*domainModel.Task, error)
		GetAllTasks(ctx context.Context) ([]domainModel.Task, error)
		GetTaskByID(ctx context.Context, id int64) (*domainModel.Task, error)
		UpdateTask(ctx context.Context, task domainModel.Task) error
		DeleteTask(ctx context.Context, id int64) error
	}
)

type service struct {
	taskRepo taskRepo
}

func New(
	taskRepo taskRepo,
) *service {
	return &service{
		taskRepo: taskRepo,
	}
}

var _ Service = (*service)(nil)

func (s service) CreateTask(ctx context.Context, task domainModel.Task) (*domainModel.Task, error) {
	return s.taskRepo.CreateTask(ctx, task)
}

func (s service) GetAllTasks(ctx context.Context) ([]domainModel.Task, error) {
	return s.taskRepo.GetAllTasks(ctx)
}

func (s service) GetTaskByID(ctx context.Context, id int64) (*domainModel.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s service) UpdateTask(ctx context.Context, task domainModel.Task) error {
	return s.taskRepo.UpdateTask(ctx, task)
}

func (s service) DeleteTask(ctx context.Context, id int64) error {
	return s.taskRepo.DeleteTask(ctx, id)
}

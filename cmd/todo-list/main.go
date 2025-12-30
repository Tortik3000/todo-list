package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tortik3000/todo-list/config"
	taskHandler "github.com/Tortik3000/todo-list/internal/api/handlers/task"
	"github.com/Tortik3000/todo-list/internal/api/middlewares"
	taskRepo "github.com/Tortik3000/todo-list/internal/repository/in_memory/task"
	taskService "github.com/Tortik3000/todo-list/internal/service/task"
)

const (
	GracefulShutdownTimeout = 5 * time.Second
)

// @title Todo List API
// @version 1.0
// @description API for managing a todo list.
// @host localhost:8080
// @BasePath /
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()

	repo := taskRepo.New()
	service := taskService.New(repo)
	handler := taskHandler.New(service)

	mux := initRouter(handler)
	handlerWithLogging := middlewares.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Server.Port),
		Handler: handlerWithLogging,
	}

	go gracefulShutdown(ctx, server)

	log.Printf("Server started, address: %s", server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("ERROR: Server: %v", err)
	}
}

func gracefulShutdown(
	ctx context.Context,
	srv *http.Server,
) {
	<-ctx.Done()
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("ERROR: Server shutdown failed: %v", err)
		return
	}
	log.Println("Shutting down todo-list")
}

func initRouter(
	taskHandler taskHandler.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		case http.MethodGet:
			taskHandler.GetAllTasks(w, r)
		default:
			methodNotAllowed(w)
		}
	})

	mux.HandleFunc("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTaskByID(w, r)
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			methodNotAllowed(w)

		}
	})
	return mux
}

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

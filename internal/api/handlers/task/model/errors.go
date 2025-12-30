package model

import (
	"encoding/json"
	"errors"
	"net/http"

	modelErr "github.com/Tortik3000/todo-list/internal/model/error"
)

var (
	ErrInvalidJSON = errors.New("invalid json")
	ErrEmptyHeader = errors.New("header should not be empty")
	ErrInvalidID   = errors.New("invalid ID")
	ErrInternalErr = errors.New("internal err")
)

func ToHttpErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var status int

	switch {
	case errors.Is(err, modelErr.ErrTaskNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrEmptyHeader) ||
		errors.Is(err, ErrInvalidJSON) ||
		errors.Is(err, ErrInvalidID):
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
		err = ErrInternalErr
	}

	w.WriteHeader(status)

	resp := map[string]string{"error": err.Error()}
	_ = json.NewEncoder(w).Encode(resp)
}

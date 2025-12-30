package http

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func ParseID(r *http.Request) (int64, error) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	return int64(id), err
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

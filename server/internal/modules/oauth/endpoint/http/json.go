package http

import (
	"encoding/json"
	"net/http"
)

// writeJSON пишет JSON-ответ с нужным статусом
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

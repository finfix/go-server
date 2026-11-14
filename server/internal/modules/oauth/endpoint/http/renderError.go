package http

import (
	"net/http"
)

// renderError показывает ошибку, которую нельзя отдать клиенту редиректом
func renderError(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = errorTemplate.Execute(w, text)
}

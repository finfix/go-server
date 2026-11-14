package http

import (
	"net/http"
)

// renderError показывает ошибку, которую нельзя отдать клиенту редиректом
func renderError(w http.ResponseWriter, status int, text string) {
	setPageHeaders(w)
	w.WriteHeader(status)
	_ = errorTemplate.Execute(w, text)
}

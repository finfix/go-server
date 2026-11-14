package service

import (
	"strings"
)

// sameResource сравнивает URL ресурса без учета завершающего слэша: клиенты передают адрес
// сервера как ввел пользователь
func sameResource(a, b string) bool {
	return strings.TrimSuffix(a, "/") == strings.TrimSuffix(b, "/")
}

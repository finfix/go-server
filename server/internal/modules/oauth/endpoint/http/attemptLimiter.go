package http

import (
	"sync"
	"time"
)

// attemptLimiter — счетчик попыток по ключу в скользящем окне, в памяти процесса. Сервер работает
// в одной реплике, поэтому общего хранилища не нужно
type attemptLimiter struct {
	mutex    sync.Mutex             // Защита карты попыток
	window   time.Duration          // Окно подсчета
	limit    int                    // Попыток в окне
	attempts map[string][]time.Time // Моменты попыток по ключу
}

// newAttemptLimiter создает счетчик попыток
func newAttemptLimiter(limit int, window time.Duration) *attemptLimiter {
	return &attemptLimiter{window: window, limit: limit, attempts: make(map[string][]time.Time)}
}

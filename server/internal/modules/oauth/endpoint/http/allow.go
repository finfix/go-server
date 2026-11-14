package http

import "time"

// allow проверяет, что по ключу не исчерпан лимит, и учитывает попытку
func (l *attemptLimiter) allow(key string) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Оставляем попытки внутри окна; заодно чистим ключи без свежих попыток, чтобы карта не росла
	now := time.Now()
	for k, times := range l.attempts {
		if len(times) != 0 && now.Sub(times[len(times)-1]) > l.window {
			delete(l.attempts, k)
		}
	}
	recent := l.attempts[key][:0]
	for _, at := range l.attempts[key] {
		if now.Sub(at) <= l.window {
			recent = append(recent, at)
		}
	}
	if len(recent) >= l.limit {
		l.attempts[key] = recent
		return false
	}
	l.attempts[key] = append(recent, now)
	return true
}

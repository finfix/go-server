package http

import (
	"net"
	"net/http"
)

// clientIP — адрес клиента. Снаружи сервер доступен только через ingress-nginx, который кладет
// настоящий адрес в X-Real-IP; без него — адрес соединения
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

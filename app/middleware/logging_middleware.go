package middleware

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: Method=%s, URL=%s, Path=%s, RemoteAddr=%s, UserAgent=%s",
			r.Method, r.URL.String(), r.URL.Path, r.RemoteAddr, r.UserAgent())

		next.ServeHTTP(w, r)

		log.Printf("Handled request: Method=%s, URL=%s", r.Method, r.URL.String())
	})
}

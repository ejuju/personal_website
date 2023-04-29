package app

import (
	"log"
	"net/http"
	"runtime/debug"
)

func newRecoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %s\n%s\n", err, debug.Stack())
					respondErrorPage(w, http.StatusInternalServerError, "fatal error")
					return
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

package api

import (
	"github.com/karnalab/karna/core"
	"net/http"
	"time"
)

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log := r.Method + " " + r.RequestURI + " " + time.Since(start).String()

		core.LogSuccessMessage(log)
	})
}

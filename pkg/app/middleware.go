package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		if path != "/metrics" {
			w.Header().Add("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

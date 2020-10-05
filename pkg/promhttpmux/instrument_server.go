package promhttpmux

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func InstrumentHttpDuration(opt prometheus.HistogramOpts) mux.MiddlewareFunc {
	obs := promauto.NewHistogramVec(opt, []string{"path", "method"})
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			labels := prometheus.Labels{"path": path, "method": r.Method}
			timer := prometheus.NewTimer(obs.With(labels))
			next.ServeHTTP(w, r)
			timer.ObserveDuration()
		})
	}
}

// InstrumentHttpDuration implements mux.MiddlewareFunc.
func InstrumentHttpInFlight(opt prometheus.GaugeOpts) mux.MiddlewareFunc {
	obs := promauto.NewGaugeVec(opt, []string{"path", "method"})
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			labels := prometheus.Labels{"path": path, "method": r.Method}
			obs.With(labels).Inc()
			defer obs.With(labels).Desc()
			next.ServeHTTP(w, r)
		})
	}
}

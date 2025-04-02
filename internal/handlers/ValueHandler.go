package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MetricsRepository interface {
	GetAllMetrics() map[string]string
	UpGauge(name string, value float64)
	UpCounter(name string, value int64)
}

func ValueHandler(repo MetricsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType, metricName := chi.URLParam(r, "type"), chi.URLParam(r, "name")
		if metricName == "" {
			http.Error(w, "Metric name cannot be empty", http.StatusNotFound)
			return
		}
		metrics := repo.GetAllMetrics()
		if value, ok := metrics[metricName]; ok && (metricType == "gauge" || metricType == "counter") {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprint(w, value)
		} else {
			http.Error(w, fmt.Sprintf("%s metric not found", metricType), http.StatusNotFound)
		}
	}
}

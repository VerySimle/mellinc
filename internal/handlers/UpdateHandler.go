package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/VerySimle/mellinc/internal/storage"
	"github.com/go-chi/chi/v5"
)

// UpdateHandler обновляет метрику
func UpdateHandler(repo MetricsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		metricValueStr := chi.URLParam(r, "value")

		if ct := r.Header.Get("Content-Type"); ct != "" && ct != "text/plain" {
			http.Error(w, "Expected Content-Type text/plain", http.StatusBadRequest)
			return
		}
		if metricName == "" {
			http.Error(w, "Metric name cannot be empty", http.StatusNotFound)
			return
		}

		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValueStr, 64)
			if err != nil || value < 0 {
				http.Error(w, "Invalid gauge value", http.StatusBadRequest)
				return
			}
			repo.UpGauge(metricName, value)
		case "counter":
			value, err := strconv.ParseInt(metricValueStr, 10, 64)
			if err != nil || value < 0 {
				http.Error(w, "Invalid counter value", http.StatusBadRequest)
				return
			}
			repo.UpCounter(metricName, value)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		// Логирование состояния (если хранилище имеет доступные поля)
		if ms, ok := repo.(*storage.MemStorage); ok {
			log.Printf("State: Gauge: %+v, Counter: %+v", ms.Gauge, ms.Counter)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, "OK")
	}
}

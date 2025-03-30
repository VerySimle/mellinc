package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/VerySimle/mellinc/internal/storage"
	"github.com/go-chi/chi/v5"
)

// RootHandler выводит все метрики в виде HTML-страницы
func AllHandler(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := repo.GetAllMetrics()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<html><head><title>Metrics</title></head><body><h1>Все метрики</h1><ul>")
		for name, value := range metrics {
			fmt.Fprintf(w, "<li>%s: %s</li>", name, value)
		}
		fmt.Fprint(w, "</ul></body></html>")
	}
}

// ValueHandler возвращает значение метрики по типу и имени
func ValueHandler(repo storage.Repository) http.HandlerFunc {
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

// UpdateHandler обновляет метрику
func UpdateHandler(repo storage.Repository) http.HandlerFunc {
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

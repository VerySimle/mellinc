// internal/handlers/update.go
package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/VerySimle/mellinc/internal/storage"
)

func UpdateHandler(ms *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что используется метод POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Проверяем заголовок Content-Type
		if contentType := r.Header.Get("Content-Type"); contentType != "text/plain" {
			http.Error(w, "Invalid Content-Type, expected text/plain", http.StatusBadRequest)
			return
		}

		// Парсим URL: /update/<тип>/<имя>/<значение>
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// Корректный url
		if len(pathParts) != 4 || pathParts[0] != "update" {
			http.Error(w, "Invalid URL format", http.StatusNotFound)
			return
		}
		metricType := pathParts[1]
		metricName := pathParts[2]
		metricValueStr := pathParts[3]

		// имя==пустое, возвращаем ошибку
		if metricName == "" {
			http.Error(w, "Metric name cannot be empty", http.StatusNotFound)
			return
		}

		// Обработка запроса gauge or counter
		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValueStr, 64)
			if err != nil {
				http.Error(w, "Invalid gauge value", http.StatusBadRequest)
				return
			}
			if value < 0 {
				http.Error(w, "Gauge value cannot be negative", http.StatusBadRequest)
				return
			}
			ms.UpGauge(metricName, value)

		case "counter":
			value, err := strconv.ParseInt(metricValueStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid counter value", http.StatusBadRequest)
				return
			}
			if value < 0 {
				http.Error(w, "Counter value cannot be negative", http.StatusBadRequest)
				return
			}
			ms.UpCounter(metricName, value)

		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}
		log.Printf("Current storage state: Gauge: %+v, Counter: %+v", ms.Gauge, ms.Counter)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

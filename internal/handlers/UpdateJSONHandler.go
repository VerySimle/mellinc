package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// JSON‑структура для обмена метриками
type UPDATEJSONMetric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`            // "gauge" или "counter"
	Delta *int64   `json:"delta,omitempty"` // для counter
	Value *float64 `json:"value,omitempty"` // для gauge
}

// UpdateJSONHandler — POST /update/
func UpdateJSONHandler(repo MetricsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			http.Error(w, "Expected Content-Type application/json", http.StatusBadRequest)
			return
		}
		var m JSONMetric
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if m.ID == "" {
			http.Error(w, "Metric id required", http.StatusBadRequest)
			return
		}

		// Сохраняем в репозиторий
		switch m.MType {
		case "gauge":
			if m.Value == nil {
				http.Error(w, "Missing value for gauge", http.StatusBadRequest)
				return
			}
			repo.UpGauge(m.ID, *m.Value)
		case "counter":
			if m.Delta == nil {
				http.Error(w, "Missing delta for counter", http.StatusBadRequest)
				return
			}
			repo.UpCounter(m.ID, *m.Delta)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		// Формируем ответ с актуальным значением
		all := repo.GetAllMetrics()
		resp := JSONMetric{ID: m.ID, MType: m.MType}
		switch m.MType {
		case "gauge":
			// преобразуем строку обратно в float64
			f, _ := strconv.ParseFloat(all[m.ID], 64)
			resp.Value = &f
		case "counter":
			i, _ := strconv.ParseInt(all[m.ID], 10, 64)
			resp.Delta = &i
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

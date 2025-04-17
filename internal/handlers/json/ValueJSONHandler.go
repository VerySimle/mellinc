package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// JSON‑структура для обмена метриками
type JSONMetric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func ValueJSONHandler(repo MetricsRepository) http.HandlerFunc {
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

		all := repo.GetAllMetrics()
		raw, ok := all[m.ID]
		if !ok {
			http.Error(w, fmt.Sprintf("%s metric not found", m.MType), http.StatusNotFound)
			return
		}

		resp := JSONMetric{ID: m.ID, MType: m.MType}
		w.Header().Set("Content-Type", "application/json")
		switch m.MType {
		case "gauge":
			v, _ := strconv.ParseFloat(raw, 64)
			resp.Value = &v
		case "counter":
			i, _ := strconv.ParseInt(raw, 10, 64)
			resp.Delta = &i
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}

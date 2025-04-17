package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RootHandler выводит все метрики в виде HTML-страницы
func AllHandler(repo MetricsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := repo.GetAllMetrics()
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(metrics)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<html><head><title>Metrics</title></head><body><h1>Все метрики</h1><ul>")
		for name, value := range metrics {
			fmt.Fprintf(w, "<li>%s: %s</li>", name, value)
		}
		fmt.Fprint(w, "</ul></body></html>")
	}
}

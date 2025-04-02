package handlers

import (
	"fmt"
	"net/http"
)

// RootHandler выводит все метрики в виде HTML-страницы
func AllHandler(repo MetricsRepository) http.HandlerFunc {
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

package main

import (
	"net/http"
	"strconv"
	"strings"
)

/*
Принимать и хранить произвольные метрики двух типов:
Для хранения метрик объявите тип MemStorage. Рекомендуем использовать тип struct с полем-коллекцией внутри (slice или map)
*/
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
	}
}

// Тип gauge, float64 — новое значение должно замещать предыдущее.
func (ms *MemStorage) UpGauge(name string, value float64) {
	ms.gauge[name] = value
}

// Тип counter, int64 — новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу
func (ms *MemStorage) UpCounter(name string, value int64) {
	ms.counter[name] += value
}

func (ms *MemStorage) mainUpdate(w http.ResponseWriter, r *http.Request) {
	//Проверяем post
	if r.Method != http.MethodPost {
		http.Error(w, "Post", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем Content-Type
	if contentType := r.Header.Get("Content-Type"); contentType != "text/plain" {
		http.Error(w, "text/plain", http.StatusBadRequest)
		return
	}
	//Парсим
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 4 {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return
	}
	metricType := pathParts[1]
	metricName := pathParts[2]
	metricValueStr := pathParts[3]
	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
		ms.UpGauge(metricName, value)

	case "counter":
		value, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid counter value", http.StatusBadRequest)
			return
		}
		ms.UpCounter(metricName, value)

	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Важно!
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func main() {
	storage := NewMemStorage()
	testServer := http.NewServeMux()
	testServer.HandleFunc(`/update/`, storage.mainUpdate)
	err := http.ListenAndServe(":8080", testServer)
	if err != nil {
		panic(err)
	}
}

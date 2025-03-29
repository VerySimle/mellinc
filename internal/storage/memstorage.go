package storage

import (
	"fmt"
	"sync"
)

// MemStorage – конкретная реализация Repository, хранящая метрики в памяти.
type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
	mu      sync.Mutex
}

// NewMemStorage – конструктор для MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

// UpGauge – обновляет метрику gauge (заменяет значение).
func (ms *MemStorage) UpGauge(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.Gauge[name] = value
}

// UpCounter – обновляет метрику counter (добавляет к существующему).
func (ms *MemStorage) UpCounter(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.Counter[name] += value
}

// GetAllMetrics – возвращает все метрики (gauge и counter) в виде map[string]string.
func (ms *MemStorage) GetAllMetrics() map[string]string {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	out := make(map[string]string)
	// Сначала скопируем gauge
	for k, v := range ms.Gauge {
		out[k] = fmt.Sprintf("%f", v) // Превращаем float64 в строку
	}
	// Затем counter
	for k, v := range ms.Counter {
		out[k] = fmt.Sprintf("%d", v) // Превращаем int64 в строку
	}
	return out
}

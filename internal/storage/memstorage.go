package storage

import (
	"fmt"
	"sync"
)

// MemStorage – конкретная реализация Repository, хранящая метрики в памяти.
type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
	Mu      sync.Mutex
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
	ms.Mu.Lock()
	defer ms.Mu.Unlock()
	ms.Gauge[name] = value
}

// UpCounter – обновляет метрику counter (добавляет к существующему).
func (ms *MemStorage) UpCounter(name string, value int64) {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()
	ms.Counter[name] += value
}

// GetAllMetrics – возвращает все метрики (gauge и counter) в виде map[string]string.
func (ms *MemStorage) GetAllMetrics() map[string]string {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()

	out := make(map[string]string)
	// Copy gauge
	for k, v := range ms.Gauge {
		out[k] = fmt.Sprintf("%f", v)
	}
	// Copy counter
	for k, v := range ms.Counter {
		out[k] = fmt.Sprintf("%d", v)
	}
	return out
}

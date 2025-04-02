package storage

import (
	"fmt"
	"sync"
)

// MemStorage – конкретная реализация Repository, хранящая метрики в памяти.
type MemStorage struct {
	Gauge     map[string]float64
	Counter   map[string]int64
	gaugeMu   sync.RWMutex
	counterMu sync.RWMutex
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
	ms.gaugeMu.Lock()
	defer ms.gaugeMu.Unlock()
	ms.Gauge[name] = value
}

// UpCounter – обновляет метрику counter (добавляет к существующему).
func (ms *MemStorage) UpCounter(name string, value int64) {
	ms.counterMu.Lock()
	defer ms.counterMu.Unlock()
	ms.Counter[name] += value
}

// GetAllMetrics – возвращает все метрики (gauge и counter) в виде map[string]string.
func (ms *MemStorage) GetAllMetrics() map[string]string {
	out := make(map[string]string)
	ms.gaugeMu.RLock()
	// Copy gauge
	for k, v := range ms.Gauge {
		out["gauge_"+k] = fmt.Sprintf("%g", v)
	}
	ms.gaugeMu.RUnlock()

	ms.counterMu.RLock()
	// Copy counter
	for k, v := range ms.Counter {
		out["counter_"+k] = fmt.Sprintf("%d", v)
	}
	ms.counterMu.RUnlock()
	return out
}

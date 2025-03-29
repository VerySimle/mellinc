package storage

type Repository interface {
	UpGauge(name string, value float64)
	UpCounter(name string, value int64)
	GetAllMetrics() map[string]string
}

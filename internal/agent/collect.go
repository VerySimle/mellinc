package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type MetricsRepository interface {
	GetAllMetrics() map[string]string
	UpGauge(name string, value float64)
	UpCounter(name string, value int64)
}

// Agent собирает метрики и отправляет их на сервер.
type Agent struct {
	repo           MetricsRepository // <-- Интерфейс репозитория
	ServerURL      string
	PollInterval   time.Duration
	ReportInterval time.Duration
	Client         *http.Client
	PollCount      int64
}

// NewAgent – конструктор агента, куда мы «внедряем» (inject) репозиторий.
func NewAgent(repo MetricsRepository, serverURL string, pollInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		repo:           repo,
		ServerURL:      serverURL,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		Client:         &http.Client{},
		PollCount:      0,
	}
}

// collectMetrics – читает runtime.MemStats и обновляет репозиторий (gauge-метрики).
func (a *Agent) collectMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Пример явного добавления ключей:
	a.repo.UpGauge("Alloc", float64(memStats.Alloc))               // Количество байт, выделенных в куче и используемых в данный момент.
	a.repo.UpGauge("BuckHashSys", float64(memStats.BuckHashSys))   // Байты, используемые для внутренней хэш-таблицы профайлера.
	a.repo.UpGauge("Frees", float64(memStats.Frees))               // Общее количество объектов, освобожденных сборщиком мусора.
	a.repo.UpGauge("GCCPUFraction", memStats.GCCPUFraction)        // Доля времени ЦП, затраченная на сборку мусора
	a.repo.UpGauge("GCSys", float64(memStats.GCSys))               // Байты, выделенные для системных структур сборщика мусора.
	a.repo.UpGauge("HeapAlloc", float64(memStats.HeapAlloc))       // Количество байт, выделенных в куче, которые всё ещё используются.
	a.repo.UpGauge("HeapIdle", float64(memStats.HeapIdle))         // Количество байт в неиспользуемых участках кучи.
	a.repo.UpGauge("HeapInuse", float64(memStats.HeapInuse))       // Количество байт в используемых участках кучи.
	a.repo.UpGauge("HeapObjects", float64(memStats.HeapObjects))   // Количество объектов, выделенных в куче.
	a.repo.UpGauge("HeapReleased", float64(memStats.HeapReleased)) // Байты, возвращённые операционной системе из кучи.
	a.repo.UpGauge("HeapSys", float64(memStats.HeapSys))           // Общее количество байт, полученных от операционной системы для кучи.
	a.repo.UpGauge("LastGC", float64(memStats.LastGC))             // Время последнего запуска сборщика мусора
	a.repo.UpGauge("Lookups", float64(memStats.Lookups))           // Количество выполненных поисков указателей
	a.repo.UpGauge("MCacheInuse", float64(memStats.MCacheInuse))   // Байты, используемые кешем для мелких объектов
	a.repo.UpGauge("MCacheSys", float64(memStats.MCacheSys))       // Общее количество байт, выделенных для mcache
	a.repo.UpGauge("MSpanInuse", float64(memStats.MSpanInuse))     // Байты, используемые для управления блоками памяти, которые находятся в использовании
	a.repo.UpGauge("MSpanSys", float64(memStats.MSpanSys))         // Общее количество байт, выделенных для mspan.
	a.repo.UpGauge("Mallocs", float64(memStats.Mallocs))           // Общее число вызовов выделения памяти
	a.repo.UpGauge("NextGC", float64(memStats.NextGC))             // Целевой объём кучи, при достижении которого будет запущен следующий сборщик мусора
	a.repo.UpGauge("NumForcedGC", float64(memStats.NumForcedGC))   // Количество циклов сборки мусора, принудительно запущенных приложением
	a.repo.UpGauge("NumGC", float64(memStats.NumGC))               // Общее количество завершённых циклов сборки мусора
	a.repo.UpGauge("OtherSys", float64(memStats.OtherSys))         // Байты, используемые для прочих системных структур
	a.repo.UpGauge("PauseTotalNs", float64(memStats.PauseTotalNs)) // Суммарное время, затраченное на паузы сборщика мусора
	a.repo.UpGauge("StackInuse", float64(memStats.StackInuse))     // Количество байт, используемых стеками горутин
	a.repo.UpGauge("StackSys", float64(memStats.StackSys))         // Общее количество байт, выделенных для стеков горутин
	a.repo.UpGauge("Sys", float64(memStats.Sys))                   // Общее количество байт, полученных от операционной системы
	a.repo.UpGauge("TotalAlloc", float64(memStats.TotalAlloc))     // Кумулятивное количество байт, когда-либо выделенных для объектов в куче
}

// updateAdditionalMetrics – добавляет PollCount (counter) и RandomValue (gauge).
func (a *Agent) updateAdditionalMetrics() {
	a.PollCount++
	a.repo.UpCounter("PollCount", 1)              // счётчик, увеличивающийся на 1 при каждом обновлении метрики из пакета
	a.repo.UpGauge("RandomValue", rand.Float64()) // обновляемое произвольное значение
}

// sendMetric – отправляет одну метрику по HTTP POST.
func (a *Agent) sendMetric(metricType, metricName, value string) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", a.ServerURL, metricType, metricName, value)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}
	return nil
}

// reportMetrics – берёт все метрики из репозитория и отправляет их на сервер.
func (a *Agent) reportMetrics() {
	metrics := a.repo.GetAllMetrics()
	for key, value := range metrics {
		metricType := "gauge"
		if key == "PollCount" {
			metricType = "counter"
		}
		err := a.sendMetric(metricType, key, value)
		if err != nil {
			log.Printf("Error sending metric %s: %v", key, err)
		}
	}
}

// Run – основной цикл агента: каждые  "a.PollInterval" собирает метрики, каждые "a.ReportInterval" отправляет.
func (a *Agent) Run() {
	pollTicker := time.NewTicker(a.PollInterval)
	reportTicker := time.NewTicker(a.ReportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			a.collectMetrics()
			a.updateAdditionalMetrics()
		case <-reportTicker.C:
			a.reportMetrics()
		}
	}
}

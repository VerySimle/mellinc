package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VerySimle/mellinc/internal/storage"
)

func TestAgent_SendMetric(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "text/plain" {
			t.Errorf("expected Content-Type text/plain, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Создаём репозиторий и агента.
	repo := storage.NewMemStorage()
	a := NewAgent(repo, server.URL, 1*time.Second, 1*time.Second)

	// Тестируем sendMetric напрямую.
	err := a.sendMetric("gauge", "testMetric", "123")
	if err != nil {
		t.Errorf("sendMetric failed: %v", err)
	}
}

func TestAgent_ReportMetrics(t *testing.T) {
	// Создаем тестовый сервер, который собирает отправленные метрики.
	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	repo := storage.NewMemStorage()
	a := NewAgent(repo, server.URL, 1*time.Second, 1*time.Second)

	// Обновляем репозиторий вручную.
	repo.UpGauge("testGauge", 42.0)
	repo.UpCounter("PollCount", 5)

	// Вызываем reportMetrics, чтобы отправить все метрики.
	a.reportMetrics()

	if len(requests) < 2 {
		t.Errorf("expected at least 2 metric requests, got %d", len(requests))
	}
	expectedPrefix := "/update/gauge/testGauge/"
	found := false
	for _, path := range requests {
		if len(path) > len(expectedPrefix) && path[:len(expectedPrefix)] == expectedPrefix {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected at least one request with prefix %s, got %+v", expectedPrefix, requests)
	}
}

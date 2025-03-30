package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/VerySimle/mellinc/internal/storage"
	"github.com/go-chi/chi/v5"
)

func TestUpdateHandler_Gauge(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		contentType    string
		metricName     string
		metricValue    string
		expectedValue  float64
		expectedStatus int
	}{
		{
			name:           "OK gauge metric",
			url:            "/update/gauge/testMetric/123.45",
			contentType:    "text/plain",
			metricName:     "testMetric",
			metricValue:    "123.45",
			expectedValue:  123.45,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "zero gauge metric",
			url:            "/update/gauge/Metric/0",
			contentType:    "text/plain",
			metricName:     "Metric",
			metricValue:    "0",
			expectedValue:  0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "- gauge value",
			url:            "/update/gauge/CqCC/-124",
			contentType:    "text/plain",
			metricName:     "CqCC",
			metricValue:    "-124",
			expectedValue:  0,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "BAD gauge value",
			url:            "/update/gauge/Invalid/abc",
			contentType:    "text/plain",
			metricName:     "Invalid",
			metricValue:    "abc",
			expectedValue:  0,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Пусто",
			url:            "/update/gauge//123",
			contentType:    "text/plain",
			metricName:     "",
			metricValue:    "123",
			expectedValue:  0,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "BAD Content-Type",
			url:            "/update/gauge/testMetric/123.45",
			contentType:    "application/json",
			metricName:     "testMetric",
			metricValue:    "123.45",
			expectedValue:  0,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMemStorage()
			handler := UpdateHandler(ms)

			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", handler)
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d for URL %s", tt.expectedStatus, resp.StatusCode, tt.url)
			}

			if tt.expectedStatus == http.StatusOK {
				metrics := ms.GetAllMetrics()
				valStr, ok := metrics[tt.metricName]
				if !ok {
					t.Errorf("Metric %s not found in storage", tt.metricName)
				}
				val, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					t.Errorf("Error parsing value for %s: %v", tt.metricName, err)
				}
				if val != tt.expectedValue {
					t.Errorf("Expected value %f for %s, got %f", tt.expectedValue, tt.metricName, val)
				}
			}
		})
	}
}

// Отдельный тест для проверки аккумулирования (накопления) значения для метрики типа counter.
func TestUpdateHandler_Counter_Accumulation(t *testing.T) {
	ms := storage.NewMemStorage()
	handler := UpdateHandler(ms)
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handler)

	// Первое обновление: устанавливаем значение 1
	req1 := httptest.NewRequest(http.MethodPost, "/update/counter/testMetric/1", nil)
	req1.Header.Set("Content-Type", "text/plain")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("Первый запрос: ожидаем статус %d, получили %d", http.StatusOK, w1.Code)
	}

	// Второе обновление: прибавляем 2
	req2 := httptest.NewRequest(http.MethodPost, "/update/counter/testMetric/2", nil)
	req2.Header.Set("Content-Type", "text/plain")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("Второй запрос: ожидаем статус %d, получили %d", http.StatusOK, w2.Code)
	}

	// Проверяем итоговое аккумулированное значение: 1 + 2 = 3
	metrics := ms.GetAllMetrics()
	valStr, ok := metrics["testMetric"]
	if !ok {
		t.Fatalf("Метрика testMetric не найдена")
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		t.Fatalf("Ошибка парсинга значения: %v", err)
	}
	if val != 3 {
		t.Errorf("Ожидаемое аккумулированное значение 3, получено %d", val)
	}
}

// Тест для проверки обработки неверного метода
func TestUpdateHandler_InvalidMethod(t *testing.T) {
	ms := storage.NewMemStorage()
	handler := UpdateHandler(ms)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handler)
	req := httptest.NewRequest(http.MethodGet, "/update/gauge/testMetric/123.45", nil)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d (MethodNotAllowed), got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

// Тест для проверки обработки неверного типа метрики
func TestUpdateHandler_InvalidType(t *testing.T) {
	ms := storage.NewMemStorage()
	handler := UpdateHandler(ms)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handler)
	req := httptest.NewRequest(http.MethodPost, "/update/invalid/testMetric/123", nil)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d (BadRequest), got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

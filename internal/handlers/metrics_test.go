package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/VerySimle/mellinc/internal/storage"
)

func TestUpdateHandler_Gauge_TableDriven(t *testing.T) {
	tests := []struct {
		url        string
		key        string
		expected   float64
		shouldFail bool
	}{
		{"/update/gauge/testMetric/123.45", "testMetric", 123.45, false},
		{"/update/gauge/Metric/0.45", "Metric", 0.45, false},
		{"/update/gauge/Test/22", "Test", 22, false},
		{"/update/gauge/CqCC/-124", "CqCC", -124, true},
		{"/update/gauge/3rEFQ/12433", "3rEFQ", 12433, false},
		{"/upload/gauge212/2121/12433", "2121", 12433, true},
	}

	for _, tt := range tests {
		ms := storage.NewMemStorage()
		handler := UpdateHandler(ms)

		req := httptest.NewRequest(http.MethodPost, tt.url, nil)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if tt.shouldFail {
			if resp.StatusCode == http.StatusOK {
				t.Errorf("For URL %s: expected failure, but got status OK", tt.url)
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("invalid URL %s: expected status OK, got %d", tt.url, resp.StatusCode)
		}

		metrics := ms.GetAllMetrics()
		valStr, ok := metrics[tt.key]
		if !ok {
			t.Errorf("For key %s: expected key present in metrics", tt.key)
		}
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			t.Errorf("For key %s: error parsing value: %v", tt.key, err)
		}
		if val != tt.expected {
			t.Errorf("For key %s: expected value %f, got %f", tt.key, tt.expected, val)
		}
	}
}

func TestUpdateHandler_Counter(t *testing.T) {
	ms := storage.NewMemStorage()
	handler := UpdateHandler(ms)

	tests := []struct {
		url        string
		key        string
		expected   int64
		shouldFail bool
	}{
		{"/update/counter/testMetric/1", "testMetric", 1, false},
		{"/update/counter/Metric/2", "Metric", 2, false},
		{"/update/counter/Test/3", "Test", 3, false},
		{"/update/counter/CqCC/-1", "CqCC", -1, true},
		{"/update/counter/3rEFQ/4", "3rEFQ", 4, false},
		{"/update/counter/testMetric/1", "testMetric", 2, false},
		{"/update/counter/Metric/1", "Metric", 3, false},
		{"/update/counter/Test/3", "Test", 6, false},
		{"/update/counter/3rEFQ/8", "3rEFQ", 12, false},
		{"/upload/test/eee/122", "eee", 122, true},
		{"/upload/counter/qqq/8", "qqq", 8, true},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodPost, tt.url, nil)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if tt.shouldFail {
			if resp.StatusCode == http.StatusOK {
				t.Errorf("For URL %s: expected failure, but got status OK", tt.url)
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("invalid URL %s: expected status OK, got %d", tt.url, resp.StatusCode)
		}

		metrics := ms.GetAllMetrics()
		valStr, ok := metrics[tt.key]
		if !ok {
			t.Errorf("For key %s: expected key present in metrics", tt.key)
		}
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			t.Errorf("For key %s: error parsing value: %v", tt.key, err)
		}
		if val != tt.expected {
			t.Errorf("For key %s: expected value %b, got %b", tt.key, tt.expected, val)
		}
	}
}

func TestUpdateHandler_InvalidMethod(t *testing.T) {
	ms := storage.NewMemStorage()
	handler := UpdateHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/update/gauge/testMetric/123.45", nil)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	handler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status MethodNotAllowed, got %d", resp.StatusCode)
	}
}

package flagsenv

import (
	"os"
	"reflect"
	"testing"
)

// TestParseFlagsAgent проверяет, что функция ParseFlagsAgent корректно
// читает значения из окружения и флагов командной строки, применяется
// правильный приоритет (флаги переопределяют ENV), и обрабатывает ошибки.
func TestParseFlagsAgent(t *testing.T) {
	// Сохраняем текущие os.Args, чтобы вернуть их после теста
	origArgs := os.Args
	// Используем defer для восстановления os.Args вне зависимости от результата теста
	defer func() { os.Args = origArgs }()

	// Таблица тестовых сценариев
	tests := []struct {
		name      string            // Описание кейса
		args      []string          // Значения os.Args для парсера
		env       map[string]string // Переменные окружения для установки перед вызовом
		want      OptionsAgent      // Ожидаемая конфигурация
		wantError bool              // Должна ли функция возвращать ошибку
	}{
		{
			name: "defaults",
			args: []string{"cmd"},
			env:  nil,
			want: OptionsAgent{Hp: "localhost:8080", Pi: 2, Ri: 10},
		},
		{
			name: "env only",
			args: []string{"cmd"},
			env: map[string]string{
				"ADDRESS":         "envhost:1111",
				"POLL_INTERVAL":   "3",
				"REPORT_INTERVAL": "15",
			},
			want: OptionsAgent{Hp: "envhost:1111", Pi: 3, Ri: 15},
		},
		{
			name: "flags only",
			args: []string{"cmd", "-a", "flaghost:2222", "-p", "5", "-r", "8"},
			env:  nil,
			want: OptionsAgent{Hp: "flaghost:2222", Pi: 5, Ri: 8},
		},
		{
			name: "env and flags override",
			args: []string{"cmd", "-a", "flaghost:3333", "-p", "6"},
			env: map[string]string{
				"ADDRESS":         "envhost:4444",
				"POLL_INTERVAL":   "7",
				"REPORT_INTERVAL": "20",
			},
			want: OptionsAgent{Hp: "flaghost:3333", Pi: 6, Ri: 20},
		},
	}

	// Пробегаем по всем сценариям
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Устанавливаем тестовые переменные окружения
			for k, v := range tc.env {
				os.Setenv(k, v)
			}
			// Гарантируем сброс ENV после теста
			t.Cleanup(func() {
				for k := range tc.env {
					os.Unsetenv(k)
				}
			})

			// Подменяем os.Args для теста
			os.Args = tc.args

			// Вызываем функцию парсинга флагов агента
			got, err := ParseFlagsAgent()

			// Если ожидается ошибка — проверяем её наличие и выходим
			if tc.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			// Если ошибка не ожидается, но она есть — тестируем провал
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Сравниваем результат с ожидаемым (флагами ENV + CLI)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

// TestParserFlagsServer проверяет парсер для сервера аналогично парсеру агента
func TestParserFlagsServer(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	tests := []struct {
		name      string
		args      []string
		env       map[string]string
		want      OptionsServer
		wantError bool
	}{
		{
			name: "defaults", // ни ENV, ни флагов
			args: []string{"cmd"},
			env:  nil,
			want: OptionsServer{Endpoint: "localhost:8080"},
		},
		{
			name: "env only", // только ENV
			args: []string{"cmd"},
			env: map[string]string{
				"ADDRESS": "srvhost:5555",
			},
			want: OptionsServer{Endpoint: "srvhost:5555"},
		},
		{
			name: "flags only", // только флаги CLI
			args: []string{"cmd", "-a", "flagserver:6666"},
			env:  nil,
			want: OptionsServer{Endpoint: "flagserver:6666"},
		},
		{
			name: "env and flags override",
			args: []string{"cmd", "-a", "overridesrv:7777"},
			env: map[string]string{
				"ADDRESS": "srvhost:8888",
			},
			want: OptionsServer{Endpoint: "overridesrv:7777"},
		},
		{
			name:      "invalid flag", // неизвестный флаг должен вызвать ошибку
			args:      []string{"cmd", "-unknown", "x"},
			env:       nil,
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				os.Setenv(k, v)
			}
			t.Cleanup(func() {
				for k := range tc.env {
					os.Unsetenv(k)
				}
			})
			os.Args = tc.args
			got, err := ParserFlagsServer()
			if tc.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

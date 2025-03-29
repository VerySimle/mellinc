// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"github.com/VerySimle/mellinc/internal/handlers"
	"github.com/VerySimle/mellinc/internal/storage"
)

func main() {
	// Создаем новое хранилище метрик
	ms := storage.NewMemStorage()

	// Создаем новый маршрутизатор (ServeMux)
	mux := http.NewServeMux()

	// Регистрируем обработчик для пути /update/
	mux.HandleFunc("/update/", handlers.UpdateHandler(ms))

	// Запускаем HTTP-сервер на порту 8080
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/VerySimle/mellinc/internal/handlers"
	"github.com/VerySimle/mellinc/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	ms := storage.NewMemStorage()
	mux := chi.NewRouter()

	// Регистрация маршрутов
	mux.Get("/", handlers.AllHandler(ms))
	mux.Post("/update/{type}/{name}/{value}", handlers.UpdateHandler(ms))
	mux.Get("/value/{type}/{name}", handlers.ValueHandler(ms))

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", mux)
}

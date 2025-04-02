package main

import (
	"log"
	"net/http"

	"github.com/VerySimle/mellinc/internal/flagsenv"
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

	confServer, err := flagsenv.ParserFlagsServer()
	if err != nil {
		log.Fatalf("Ошибка парсинга конфигурации сервера: %v", err)
	}

	//Вывод в терминал
	log.Printf("Server started on %s", confServer.Endpoint)
	if err := http.ListenAndServe(confServer.Endpoint, mux); err != nil {
		log.Fatal(err)
	}

}

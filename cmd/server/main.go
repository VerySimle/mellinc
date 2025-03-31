package main

import (
	"flag"
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

	//Флаг для изменения порта сервера
	endpoint := flag.String("a", "localhost:8080", "input Port")
	flag.Parse()
	addr := *endpoint

	//Вывод в терминал :endpoint
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}

}

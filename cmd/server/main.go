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

	//Флаг для изменения порта сервера
	/*
		endpoint := flag.String("a", "localhost:8080", "input Port")
		flag.Parse()
	*/

	flagsenv.ParserFlagsServer()

	//addr := *endpoint

	//Вывод в терминал
	log.Printf("Server started on %s", flagsenv.ConfServer.Endpoint)
	if err := http.ListenAndServe(flagsenv.ConfServer.Endpoint, mux); err != nil {
		log.Fatal(err)
	}

}

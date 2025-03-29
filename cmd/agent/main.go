package main

import (
	"log"
	"time"

	"github.com/VerySimle/mellinc/internal/agent"
	"github.com/VerySimle/mellinc/internal/storage"
)

func main() {
	// Создаём репозиторий (хранилище в памяти)
	repo := storage.NewMemStorage()

	// Создаём агента, передавая репозиторий и нужные интервалы
	a := agent.NewAgent(repo, "http://localhost:8080", 2*time.Second, 10*time.Second)

	log.Println("Agent started")
	a.Run()
}

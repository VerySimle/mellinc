package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/VerySimle/mellinc/internal/agent"
	"github.com/VerySimle/mellinc/internal/storage"
)

func main() {
	// Создаём репозиторий (хранилище в памяти)
	repo := storage.NewMemStorage()
	var options struct {
		hp string
		pi int
		ri int
	}

	flag.StringVar(&options.hp, "a", "localhost:8080", "Адрес и порт хоста")
	flag.IntVar(&options.pi, "p", 2, "pollInterval")
	flag.IntVar(&options.ri, "r", 4, "reportInterval")
	flag.Parse()
	addr := fmt.Sprintf("http://%s", options.hp)
	// Создаём агента, передавая репозиторий и нужные интервалы
	a := agent.NewAgent(repo, addr, time.Duration(options.pi)*time.Second, time.Duration(options.ri)*time.Second)

	log.Println("Agent started")
	a.Run()

}

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/VerySimle/mellinc/internal/agent"
	"github.com/VerySimle/mellinc/internal/flagsenv"
	"github.com/VerySimle/mellinc/internal/storage"
)

func main() {
	// Создаём репозиторий (хранилище в памяти)
	repo := storage.NewMemStorage()
	// Регистрируем и парсим флаги из пакета flagsenv
	flagsenv.ParseFlagsAgent()

	addr := fmt.Sprintf("http://%s", flagsenv.ConfAgent.Hp)
	// Создаём агента, передавая репозиторий и нужные интервалы
	a := agent.NewAgent(repo, addr, time.Duration(flagsenv.ConfAgent.Pi)*time.Second, time.Duration(flagsenv.ConfAgent.Ri)*time.Second)

	log.Println("Agent started")
	a.Run()

}

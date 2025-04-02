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
	confAgent, err := flagsenv.ParseFlagsAgent()
	if err != nil {
		log.Fatalf("Ошибка парсинга конфигурации агента: %v", err)
	}

	addr := fmt.Sprintf("http://%s", confAgent.Hp)
	// Создаём агента, передавая репозиторий и нужные интервалы
	a := agent.NewAgent(repo, addr, time.Duration(confAgent.Pi)*time.Second, time.Duration(confAgent.Ri)*time.Second)

	log.Printf("Agent started host - %s, POLL_INTERVAL - %ds, REPORT_INTERVAL - %ds", addr, confAgent.Pi, confAgent.Ri)
	a.Run()

}

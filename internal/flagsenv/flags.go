package flagsenv

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type OptionsAgent struct {
	Hp string `env:"ADDRESS" envDefault:"localhost:8080"`
	Pi int    `env:"POLL_INTERVAL" envDefault:"2"`
	Ri int    `env:"REPORT_INTERVAL" envDefault:"10"`
}

type OptionsServer struct {
	Endpoint string `env:"ADDRESS" envDefault:"localhost:8080"`
}

var ConfAgent OptionsAgent
var ConfServer OptionsServer

func ParseFlagsAgent() (OptionsAgent, error) {
	var confAgent OptionsAgent
	if err := env.Parse(&confAgent); err != nil {
		fmt.Printf("Ошибка парсинга: %+v\n", err)
		return OptionsAgent{}, err
	}
	fmt.Printf("Конфигурация агента: %+v\n", confAgent)
	return confAgent, nil
}

func ParserFlagsServer() (OptionsServer, error) {
	var confServer OptionsServer
	if err := env.Parse(&confServer); err != nil {
		fmt.Printf("Ошибка парсинга: %+v\n", err)
		return OptionsServer{}, err
	}
	fmt.Printf("Конфигурация сервера: %+v\n", confServer)
	return confServer, nil
}

package flagsenv

import (
	"flag"
	"fmt"
	"os"

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

func ParseFlagsAgent() (OptionsAgent, error) {
	var cfg OptionsAgent
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Ошибка парсинга: %+v\n", err)
		return OptionsAgent{}, err
	}
	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.StringVar(&cfg.Hp, "a", cfg.Hp, "Адрес и порт хоста")
	fs.IntVar(&cfg.Pi, "p", cfg.Pi, "Интервал опроса")
	fs.IntVar(&cfg.Ri, "r", cfg.Ri, "Интервал отчётов")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return OptionsAgent{}, err
	}
	return cfg, nil
}

func ParserFlagsServer() (OptionsServer, error) {
	var cfg OptionsServer
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Ошибка парсинга: %+v\n", err)
		return OptionsServer{}, err
	}
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Endpoint, "a", cfg.Endpoint, "input Port")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return OptionsServer{}, err
	}
	return cfg, nil
}

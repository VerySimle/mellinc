package flagsenv

import (
	"flag"
	"os"
	"strconv"
)

type OptionsAgent struct {
	Hp string
	Pi int
	Ri int
}

type OptionsServer struct {
	Endpoint string
}

var ConfAgent OptionsAgent
var ConfServer OptionsServer

func ParseFlagsAgent() {
	flag.StringVar(&ConfAgent.Hp, "a", "localhost:8080", "Адрес и порт хоста")
	flag.IntVar(&ConfAgent.Pi, "p", 2, "Интервал опроса")
	flag.IntVar(&ConfAgent.Ri, "r", 4, "Интервал отчётов")
	flag.Parse()
	if env := os.Getenv("ADDRESS"); env != "" {
		ConfAgent.Hp = env
	}
	if env := os.Getenv("REPORT_INTERVAL"); env != "" {
		if value, err := strconv.Atoi(env); err == nil {
			ConfAgent.Ri = value
		}
	}
	if env := os.Getenv("POLL_INTERVAL"); env != "" {
		if value, err := strconv.Atoi(env); err == nil {
			ConfAgent.Pi = value
		}
	}

}

func ParserFlagsServer() {
	flag.StringVar(&ConfServer.Endpoint, "a", "localhost:8080", "input Port")
	flag.Parse()
	if env := os.Getenv("ADDRESS"); env != "" {
		ConfServer.Endpoint = env
	}
}

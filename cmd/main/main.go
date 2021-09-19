package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/242617/notify/config"
	"github.com/242617/notify/server"
	"github.com/242617/notify/telegram"
)

const (
	TimeFormat  = "15:04:05"
	ServiceName = "notify"
)

var serverAddress string

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		TimeFormat: TimeFormat,
		Out:        os.Stderr,
	})
}

func main() {
	flag.StringVar(&serverAddress, "address", "0.0.0.0:8080", "Server address")
	flag.Parse()

	b, err := os.ReadFile("/etc/service/config.yaml")
	if err != nil {
		log.Debug().Err(err).Msg("cannot open configuration file")
		select {}
	}

	var cfg config.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal().Err(err).Msg("cannot unmarshal config")
	}

	telegram, err := telegram.New(cfg.Telegram)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create telegram service")
	}

	server, err := server.New(serverAddress, telegram)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	telegram.Start(context.TODO())
	go func() {
		if err := server.Start(context.TODO()); err != nil {
			log.Fatal().Err(err).Msg("cannot start server")
		}
	}()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	log.Info().Str("service", ServiceName).Msg("stop")
	telegram.Stop(context.TODO())
	server.Stop(context.TODO())
}

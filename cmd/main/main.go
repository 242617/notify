package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
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

var (
	serverAddress string
	consulAddress string
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		TimeFormat: TimeFormat,
		Out:        os.Stderr,
	})
}

func main() {
	flag.StringVar(&serverAddress, "address", "0.0.0.0:8080", "Server address")
	flag.StringVar(&consulAddress, "consul", "0.0.0.0:8500", "Consul address")
	flag.Parse()

	if os.Getenv("CONSUL_HTTP_TOKEN") == "" {
		log.Fatal().Msg("empty token")
	}

	client, err := api.NewClient(&api.Config{Address: consulAddress})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create client")
	}

	service, err := connect.NewService(ServiceName, client)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new service")
	}
	defer service.Close()

	<-service.ReadyWait()

	log.Info().Str("service", ServiceName).Msg("start")

	pair, _, err := client.KV().Get(ServiceName, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get kv")
	}

	var cfg config.Config
	if err := yaml.Unmarshal(pair.Value, &cfg); err != nil {
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

	//

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

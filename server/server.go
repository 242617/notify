package server

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/242617/notify/telegram"
)

func New(address string, telegram telegram.API) (*server, error) {
	s := server{
		srv:      grpc.NewServer(),
		address:  address,
		telegram: telegram,
	}
	RegisterNotifyServer(s.srv, &s)
	grpc_health_v1.RegisterHealthServer(s.srv, &s)
	reflection.Register(s.srv)
	return &s, nil
}

type server struct {
	srv      *grpc.Server
	address  string
	telegram telegram.API
}

func (s *server) Start(context.Context) error {
	conn, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Error().Err(err).Msg("cannot listen")
		return err
	}
	log.Debug().Str("s.address", s.address).Msg("start listening")
	return s.srv.Serve(conn)
}

func (s *server) Stop(context.Context) error {
	s.srv.GracefulStop()
	log.Debug().Str("s.address", s.address).Msg("stop listening")
	return nil
}

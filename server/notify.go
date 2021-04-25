package server

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (s *server) Notify(ctx context.Context, in *NotifyRequest) (*NotifyResponse, error) {
	log.Debug().Str("in.Message", in.Message).Msg("Notify")
	if err := s.telegram.Notify(ctx, in.Message); err != nil {
		log.Error().Err(err).Str("in.Message", in.Message).Msg("cannot notify")
		return nil, err
	}
	return &NotifyResponse{}, nil
}

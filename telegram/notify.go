package telegram

import (
	"context"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (t *telegram) Notify(_ context.Context, msg string) error {
	log.Debug().Str("msg", msg).Msg("Notify")
	for _, id := range t.cfg.Recipients {
		_, err := t.bot.Send(&tb.User{ID: id}, msg)
		log.Debug().Int("id", id).Str("msg", msg).Msg("Notify")
		if err != nil {
			log.Error().Err(err).Str("msg", msg).Msg("cannot send message")
			return err
		}
	}
	return nil
}

package telegram

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

type API interface {
	Notify(ctx context.Context, msg string) error
}

func New(cfg Config) (*telegram, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot create new bot")
		return nil, err
	}

	bot.Handle(tb.OnText, func(msg *tb.Message) {
		bot.Send(msg.Sender, "Hello!")
	})

	return &telegram{
		cfg: cfg,
		bot: bot,
	}, nil
}

type telegram struct {
	cfg Config
	bot *tb.Bot
}

func (t *telegram) Start(context.Context) error {
	go t.bot.Start()
	return nil
}

func (t *telegram) Stop(context.Context) error {
	go t.bot.Stop()
	return nil
}

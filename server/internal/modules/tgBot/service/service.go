package service

import (
	"go.opentelemetry.io/otel"
	"gopkg.in/telebot.v3"
	"pkg/log"
	"server/internal/config"

	"server/internal/utils/errors"
)

var tracer = otel.Tracer("/server/internal/modules/tgBot/service")

type TgBotService struct {
	Bot  *telebot.Bot
	chat *telebot.Chat
}

func NewTgBotService() (*TgBotService, error) {

	bot, err := telebot.NewBot(telebot.Settings{
		URL:     "",
		Token:   config.Load().Telegram.Token,
		Updates: 0,
		Poller: &telebot.LongPoller{
			Limit:          0,
			Timeout:        0,
			LastUpdateID:   0,
			AllowedUpdates: nil,
		},
		Synchronous: false,
		Verbose:     false,
		ParseMode:   telebot.ModeHTML,
		OnError: func(err error, c telebot.Context) {
			log.Error(err)
		},
		Client:  nil,
		Offline: !config.Load().Telegram.IsEnabled,
	})
	if err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	var chat *telebot.Chat
	if config.Load().Telegram.IsEnabled {
		if chat, err = bot.ChatByID(config.Load().Telegram.ChatID); err != nil {
			return nil, errors.InternalServer.Wrap(err)
		}
	}

	return &TgBotService{
		Bot:  bot,
		chat: chat,
	}, nil
}

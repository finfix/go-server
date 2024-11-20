package service

import (
	"context"

	"gopkg.in/telebot.v3"

	"pkg/errors"
	"pkg/log"
)

type TgBotService struct {
	Bot  *telebot.Bot
	Chat *telebot.Chat

	isOn bool
}

func NewTgBotService(
	token string,
	chatID int64,
	isOn bool,
) (*TgBotService, error) {

	if !isOn {
		log.Warning(context.Background(), "Telegram bot is off", log.SkipThisCallOption())
		return &TgBotService{
			Bot:  nil,
			Chat: nil,
			isOn: isOn,
		}, nil
	}

	bot, err := telebot.NewBot(telebot.Settings{
		URL:         "",
		Token:       token,
		Updates:     0,
		Poller:      nil,
		Synchronous: false,
		Verbose:     false,
		ParseMode:   telebot.ModeHTML,
		OnError:     nil,
		Client:      nil,
		Offline:     false,
	})
	if err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	chat, err := bot.ChatByID(chatID)
	if err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	return &TgBotService{
		Bot:  bot,
		Chat: chat,
		isOn: isOn,
	}, nil
}

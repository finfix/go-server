package service

import (
	"context"
	"server/internal/config"

	"server/internal/utils/errors"

	"server/internal/modules/tgBot/model"
)

// SendMessage отправляет сообщение пользователю в телеграм
func (s *TgBotService) SendMessage(ctx context.Context, req model.SendMessageReq) error {
	ctx, span := tracer.Start(ctx, "SendMessage")
	defer span.End()

	if !config.Load().Telegram.IsEnabled {
		return nil
	}

	if _, err := s.Bot.Send(s.chat, req.Message); err != nil {
		return errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}

	return nil
}

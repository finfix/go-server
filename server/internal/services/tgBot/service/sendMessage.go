package service

import (
	"context"

	"pkg/errors"
	"pkg/log"

	"server/internal/services/tgBot/model"
)

// SendMessage отправляет сообщение пользователю в телеграм
func (s *TgBotService) SendMessage(ctx context.Context, req model.SendMessageReq) error {
	ctx, span := tracer.Start(ctx, "SendMessage")
	defer span.End()

	if !s.isOn {
		log.Warning(ctx, "Вызвана функция SendMessage. Пуши выключены")
		return nil
	}

	if _, err := s.Bot.Send(s.Chat, req.Message); err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	return nil
}

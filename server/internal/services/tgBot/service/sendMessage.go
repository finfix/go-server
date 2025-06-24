package service

import (
	"context"

	"server/internal/utils/errors"

	"server/internal/services/tgBot/model"
)

// SendMessage отправляет сообщение пользователю в телеграм
func (s *TgBotService) SendMessage(ctx context.Context, req model.SendMessageReq) error {
	ctx, span := tracer.Start(ctx, "SendMessage")
	defer span.End()

	if _, err := s.Bot.Send(s.chat, req.Message); err != nil {
		return errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}

	return nil
}

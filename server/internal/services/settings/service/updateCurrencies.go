package service

import (
	"context"
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"pkg/log"

	settingsModel "server/internal/services/settings/model"
	"server/internal/services/settings/network"
	"server/internal/services/settings/service/utils"
	"server/internal/services/tgBot/model"
)

// UpdateCurrencies обновляет курсы валют
func (s *SettingsService) UpdateCurrencies(ctx context.Context, req settingsModel.UpdateCurrenciesReq) error {
	ctx, span := tracer.Start(ctx, "UpdateCurrencies")
	defer span.End()

	// Проверяем, что пользователь администратор
	err := s.checkAdmin(ctx, req.Necessary.UserID)
	if err != nil {
		return err
	}

	const updateCurrenciesTemplate = "<b>📈 Курс валют успешно обновлен</b>\n\nUSD: %v₽\nBTC: %v$"

	var tgMessage model.SendMessageReq

	defer func() {
		err := s.tgBot.SendMessage(ctx, tgMessage)
		if err != nil {
			log.WithContextParams(ctx).Error(err)
		}
	}()

	// Получаем курсы валют от провайдера данных
	rates, err := network.GetCurrencyRates(ctx, s.credentials.CurrencyProviderAPIKey)
	if err != nil {
		tgMessage.Message += fmt.Sprintf("Не смогли получить курсы валют от провайдера\n\n%v", err.Error())
		return err
	}
	tgMessage.Message += "Успешно получили курсы валют от провайдера\n"

	// Обновляем курсы валют в БД
	if err = s.settingsRepository.UpdateCurrencies(ctx, rates); err != nil {
		tgMessage.Message += fmt.Sprintf("Не смогли обновить курсы валют в базе данных\n\n%v", err.Error())
		return err
	}

	p := message.NewPrinter(language.Russian)

	usdrubRate, _ := utils.GetRate(rates, "USD", "RUB").Float64()
	btcusdRate, _ := utils.GetRate(rates, "BTC", "USD").Float64()

	tgMessage.Message = fmt.Sprintf(
		updateCurrenciesTemplate,
		p.Sprintf("%.2f", usdrubRate),
		p.Sprintf("%.0f", btcusdRate),
	)

	return nil
}

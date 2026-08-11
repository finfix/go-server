package service

import (
	"context"
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"pkg/log"
	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	settingsModel "server/internal/modules/settings/model"
	"server/internal/modules/settings/network"
	"server/internal/modules/settings/service/utils"
	"server/internal/modules/tgBot/model"
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

	// Получаем текущие курсы валют для слепка "до" в аудит-логе
	currenciesBefore, err := s.settingsRepository.GetCurrencies(ctx)
	if err != nil {
		return err
	}
	currenciesBeforeMap := slices.ToMap(currenciesBefore, func(c settingsModel.Currency) string { return c.Slug })

	err = s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Обновляем курсы валют в БД
		if err := s.settingsRepository.UpdateCurrencies(ctxTx, rates); err != nil {
			return err
		}

		// Получаем актуальные курсы валют из БД для слепков "после" в аудит-логе
		currenciesAfter, err := s.settingsRepository.GetCurrencies(ctxTx)
		if err != nil {
			return err
		}
		currenciesAfterMap := slices.ToMap(currenciesAfter, func(c settingsModel.Currency) string { return c.Slug })

		// Фиксируем изменение каждого курса валюты в аудит-логе
		for slug := range rates {
			method := auditLogMethod.Update
			var before any
			if currencyBefore, ok := currenciesBeforeMap[slug]; ok {
				before = currencyBefore
			} else {
				method = auditLogMethod.Create
			}

			if _, err := s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
				Entity:   auditLogEntity.Currency,
				Method:   method,
				EntityID: slug,
				Before:   before,
				After:    currenciesAfterMap[slug],
				UserID:   req.Necessary.UserID,
				DeviceID: req.Necessary.DeviceID,
			}); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
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

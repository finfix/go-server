package scheduler

import (
	"context"
	"time"

	"pkg/errors"
	"pkg/log"
	"server/internal/utils/necessary"

	"server/internal/services/settings/model"
)

func (s *Scheduler) Start(ctx context.Context) error {

	// Обновление валют
	_, err := s.cron.AddFunc("@daily", func() { // Every day at 00:00 UTC
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		ctx, span := tracer.Start(ctx, "UpdateCurrencies")
		defer span.End()

		if err := s.settingsService.UpdateCurrencies(ctx, model.UpdateCurrenciesReq{
			Necessary: necessary.NecessaryUserInformation{
				UserID:   adminUser,
				DeviceID: "system",
			},
		}); err != nil {
			log.Error(context.Background(), err)
		}
	})
	if err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	s.cron.Start()

	return nil
}

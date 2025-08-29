package scheduler

import (
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"

	settingsService "server/internal/modules/settings/service"
)

const adminUser = 15

var tracer = otel.Tracer("/server/internal/modules/scheduler/service")

type Scheduler struct {
	cron            *cron.Cron
	settingsService *settingsService.SettingsService
}

func NewScheduler(
	settingsService *settingsService.SettingsService,

) *Scheduler {
	return &Scheduler{
		cron:            cron.New(),
		settingsService: settingsService,
	}
}

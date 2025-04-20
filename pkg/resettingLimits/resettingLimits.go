package resettingLimits

import (
	"context"
	"github.com/ether-echo/user-service/pkg/logger"
	"github.com/robfig/cron/v3"
)

var (
	log = logger.Logger().Named("resettingLimits").Sugar()
)

type IRepository interface {
	ResetFlags(ctx context.Context) error
}

type ResettingLimits struct {
	cron        *cron.Cron
	IRepository IRepository
}

func NewResettingLimits(IRepository IRepository) *ResettingLimits {
	return &ResettingLimits{
		cron:        cron.New(),
		IRepository: IRepository,
	}
}

func (r *ResettingLimits) ResetLimit(ctx context.Context) {
	_, err := r.cron.AddFunc("0 0 * * *", func() {
		err := r.IRepository.ResetFlags(ctx)
		if err != nil {
			log.Errorf("Resetting flags failed: %v", err)
		} else {
			log.Infof("Resetting flags succeeded")
		}
	})
	if err != nil {
		log.Errorf("Resetting flags failed: %v", err)
	}

	r.cron.Start()

	log.Info("Resetting flags started")

	select {}
}

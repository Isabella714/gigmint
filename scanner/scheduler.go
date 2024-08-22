package scanner

import (
	"time"

	"github.com/go-co-op/gocron"
)

func StartJobScheduler() error {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonModeAll()

	_, err := RegisterBlockEventScanner(scheduler)
	if err != nil {
		return err
	}

	scheduler.StartAsync()
	return nil
}

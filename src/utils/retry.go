package utils

import (
	"context"
	"errors"
	"log"
	"time"
)

const initialSleepInterval = 500
const maxSleepInterval = 16 * 1000

func nextSleepInterval(ms time.Duration) time.Duration {
	if ms == 0 {
		return initialSleepInterval
	}
	if ms >= maxSleepInterval {
		return ms
	}
	return ms * 2
}

func Retry(task func() error, desc string, healthCallback func(bool)) (err error) {
	var sleepInterval time.Duration = 0
	for {
		err = task()
		if err == nil {
			log.Printf("%v succeeded", desc)
			healthCallback(true)
			return err
		}
		if errors.Is(err, context.Canceled) {
			return err
		}

		log.Printf("%v failed with: %v, would retry", desc, err)
		healthCallback(false)

		time.Sleep(sleepInterval * time.Millisecond)
		sleepInterval = nextSleepInterval(sleepInterval)
	}
}

package utils

import (
	"context"
	"errors"
	"log"
	"time"
)

func Retry(task func() error, desc string, healthCallback func(bool)) (err error) {
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
		time.Sleep(10 * time.Second)
	}
}

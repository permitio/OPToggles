package targets

import (
	"context"
	"log"
	"optoggles/trackers"
)

type LogPublisher struct {
}

func NewLogPublisher() *LogPublisher {
	return &LogPublisher{}
}

func (pp *LogPublisher) Work(ctx context.Context, flagsChan trackers.ToggleEvents) error {
	for {
		select {
		case queryResult := <-flagsChan:
			log.Printf("New configuration for user-authorized toggles")
			log.Printf("* %s would be enabled for users: %s", queryResult.Toggle.TargetSpec.Key, queryResult.Users)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

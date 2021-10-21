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

func (pp *LogPublisher) SyncForever(ctx context.Context, flagsChan trackers.ToggleUpdates) error {
	for {
		select {
		case ToggleUpdate := <-flagsChan:
			log.Printf("New configuration for user-authorized toggles")
			log.Printf("* %s would be enabled for users: %s", ToggleUpdate.Toggle.Spec.Key, ToggleUpdate.Users)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

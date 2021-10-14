package targets

import (
	"context"
	"log"
	"optoggles/trackers"
)

type LogPublisher struct {
	flagsChan trackers.ToggleEvents
}

func NewLogPublisher(events trackers.ToggleEvents) *LogPublisher {
	return &LogPublisher{
		// This channel is not buffered, writer would block until the last operation is finished
		flagsChan: events,
	}
}

func (pp *LogPublisher) Work(ctx context.Context) error {
	for {
		select {
		case queryResult := <-pp.flagsChan:
			log.Printf("New configuration for user-authorized toggles")
			log.Printf("* %s would be enabled for users: %s", queryResult.Toggle.TargetSpec.Key, queryResult.Users)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
package ldpublisher

import (
	"context"
	"log"
)

type PrintPublisher struct {
	flagsChan chan map[string][]string
}

func NewPrintPublisher() (*PrintPublisher, error) {
	return &PrintPublisher{
		// This channel is not buffered, writer would block until the last operation is finished
		flagsChan: make(chan map[string][]string),
	}, nil
}

func (pp *PrintPublisher) Publish(usersToFlags map[string][]string) {
	pp.flagsChan <- usersToFlags
}

func (pp *PrintPublisher) Work(ctx context.Context) error {
	for {
		select {
		case usersToFlags := <-pp.flagsChan:
			log.Printf("New configuration for user-authorized toggles")
			for flag, users := range usersToFlags {
				log.Printf("* %s would be enabled for users: %s", flag, users)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
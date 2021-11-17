package targets

import (
	"context"
	"errors"
	"fmt"
	"log"
	"optoggles/config"
	"optoggles/httpserver"
	"optoggles/trackers"
	"optoggles/utils"
)

type Target interface {
	CreateToggle(ctx context.Context, key string, spec map[string]interface{}) error
	UpdateToggleWithUsers(ctx context.Context, key string, users []string) error
}

type TargetSync struct {
	target Target
}

func (ts *TargetSync) InitTarget(ctx context.Context, toggles []config.ToggleConfig) error {
	for _, toggle := range toggles {
		if err := ts.target.CreateToggle(ctx, toggle.Key, toggle.Spec); err != nil {
			return err
		}
	}
	return nil
}

func (ts *TargetSync) SyncForever(ctx context.Context, flagsChan trackers.ToggleUpdates) (err error) {
	for {
		select {
		case updates := <-flagsChan:
			for _, update := range updates {
				log.Printf("toggle %s would be enabled for users: %s", update.Toggle.Key, update.Users)
				err = utils.Retry(func() error {
					return ts.target.UpdateToggleWithUsers(ctx, update.Toggle.Key, update.Users)
				}, fmt.Sprintf("updating %v's users", update.Toggle.Key),
					httpserver.HealthServ.SetTargetSyncHealth)
				if err != nil {
					return err
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func NewTargetSync(config config.TargetConfig) (ts *TargetSync, err error) {
	ts = &TargetSync{}

	switch config.TargetType {
	case "launchdarkly":
		ts.target, err = NewLaunchdarklyTarget(config.TargetSpec)
	case "restapi":
		ts.target, err = NewRestApiTarget(config.TargetSpec)
	default:
		err = errors.New("target type isn't supported. currently supports: [launchdarkly, restapi]")
	}
	return
}

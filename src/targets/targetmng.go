package targets

import (
	"context"
	"errors"
	"log"
	"optoggles/config"
	"optoggles/trackers"
)

type Target interface {
	CreateToggle(ctx context.Context, key string, spec map[string]interface{}) error
	UpdateToggleWithUsers(ctx context.Context, key string, users []string) error
}

func targetFactory(config config.TargetConfig) (target Target, err error) {
	switch config.TargetType {
	case "launchdarkly":
		target, err = NewLaunchdarklyTarget(config.TargetSpec)
	case "http":
		target, err = NewHttpTarget(config.TargetSpec)
	default:
		err = errors.New("target type isn't supported. currently only supports 'launchdarkly'")
	}
	return
}

func InitTarget(ctx context.Context, config config.TargetConfig, toggles []config.ToggleConfig) (target Target, err error) {
	target, err = targetFactory(config)
	if err == nil {
		for _, toggle := range toggles {
			if err = target.CreateToggle(ctx, toggle.Key, toggle.Spec); err != nil {
				return
			}
		}
	}
	return
}

func SyncTargetForever(ctx context.Context, target Target, flagsChan trackers.ToggleUpdates) error {
	for {
		select {
		case ToggleUpdate := <-flagsChan:
			log.Printf("* %s would be enabled for users: %s", ToggleUpdate.Toggle.Key, ToggleUpdate.Users)
			if err := target.UpdateToggleWithUsers(ctx, ToggleUpdate.Toggle.Key, ToggleUpdate.Users); err != nil {
				// TODO: return error? retries?
				log.Printf("updating flag's target users failed: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

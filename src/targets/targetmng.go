package targets

import (
	"context"
	"errors"
	"optoggles/config"
	"optoggles/trackers"
)

type Target interface {
	CreateToggle(ctx context.Context, toggle config.ToggleConfig) error
	SyncForever(ctx context.Context, flagsChan trackers.ToggleUpdates) error
}

func InitTarget(ctx context.Context, config config.TargetConfig, toggles []config.ToggleConfig) (target Target, err error) {
	if config.TargetType != "launchdarkly" {
		err = errors.New("target type isn't supported. currently only supports 'launchdarkly'")
		return
	}

	target = NewLaunchdarklyTarget(config.LaunchdarklyToken)

	for _, toggle := range toggles {
		if err = target.CreateToggle(ctx, toggle); err != nil {
			return
		}
	}
	return
}

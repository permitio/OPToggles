package targets

import (
	"context"
	"errors"
	"optoggles/config"
	"optoggles/trackers"
)

type Target interface {
	CreateFlag(ctx context.Context, toggle config.ToggleConfig) error
	Work(ctx context.Context, flagsChan trackers.ToggleEvents) error
}

func InitTarget(ctx context.Context, config config.TargetConfig, toggles []config.ToggleConfig) (target Target, err error) {
	if config.TargetType != "launchdarkly" {
		err = errors.New("target type isn't supported. currently only supports 'launchdarkly'")
		return
	}

	target = NewLaunchdarklyTarget(config.LaunchdarklyToken)

	for _, toggle := range toggles {
		if err = target.CreateFlag(ctx, toggle); err != nil {
			return
		}
	}
	return
}

package trackers

import (
	"context"
	"optoggles/config"
)

type PolicyTracker interface {
	Track(ctx context.Context)
	Close()
}

type ToggleUpdate struct {
	Toggle config.ToggleConfig
	Users []string
}
type ToggleUpdates chan ToggleUpdate

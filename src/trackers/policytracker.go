package trackers

import (
	"context"
	"optoggles/config"
)

type PolicyTracker interface {
	Track(ctx context.Context)
	GetChanges(ctx context.Context)
	Close()
}

type QueryResult struct {
	Toggle config.ToggleConfig
	Users []string
}

type ToggleEvents chan QueryResult
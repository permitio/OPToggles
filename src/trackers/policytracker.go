package trackers

import (
	"context"
	"errors"
	"fmt"
	"optoggles/config"
)

type ToggleQueryResult struct {
	Toggle config.ToggleConfig
	Users  []string
}
type ToggleUpdates chan []ToggleQueryResult

type PolicyTracker struct {
	togglesBySource map[string][]config.ToggleConfig
	opaTrackers     []*OpaTracker
	Updates         ToggleUpdates
}

func (pt *PolicyTracker) InitTrackers(ctx context.Context) error {
	for _, opaTracker := range pt.opaTrackers {
		if err := opaTracker.Init(ctx); err != nil {
			return err
		}

		toggles := pt.togglesBySource[opaTracker.opal.config.Id]
		if err := opaTracker.QueryOnce(ctx, toggles, pt.Updates); err != nil {
			return err
		}
	}
	return nil
}

func (pt *PolicyTracker) TrackAll(ctx context.Context) {
	for _, opaTracker := range pt.opaTrackers {
		toggles := pt.togglesBySource[opaTracker.opal.config.Id]
		go opaTracker.Track(ctx, toggles, pt.Updates)
	}
}

func NewPolicyTracker(sources []config.OpalConfig, toggles []config.ToggleConfig) (pt *PolicyTracker, err error) {
	pt = &PolicyTracker{}
	pt.Updates = make(ToggleUpdates, len(sources)) // Have buffer for one update per source

	pt.togglesBySource = make(map[string][]config.ToggleConfig)
	for _, toggle := range toggles {
		srcId := toggle.UsersDocument.Source

		if _, sourceExists := pt.togglesBySource[srcId]; !sourceExists {
			err = errors.New(fmt.Sprintf("toggle %v points to non existing source %v", toggle.Key, srcId))
			for _, src := range sources {
				if src.Id == srcId {
					err = nil
					break
				}
			}
			if err != nil {
				return nil, err
			}
		}
		pt.togglesBySource[srcId] = append(pt.togglesBySource[srcId], toggle)
	}

	pt.opaTrackers = []*OpaTracker{}
	for _, source := range sources {
		tracker := NewOpaTracker(source)
		pt.opaTrackers = append(pt.opaTrackers, tracker)
	}

	return
}

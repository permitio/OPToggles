package trackers

import (
	"context"
	"log"
	"optoggles/config"
)

type OpaTracker struct {
	opal       *OpalClient
	closeTrack context.CancelFunc // TODO: is that needed?
}

func (ot *OpaTracker) Track(ctx context.Context, toggles []config.ToggleConfig, updates ToggleUpdates) (err error) {
	// If one of the tracker's goroutine ends, finish all goroutrines
	ctx, ot.closeTrack = context.WithCancel(ctx)

	var opaClient *OpaClient
	if opaClient, err = ot.opal.getOpaClient(ctx); err != nil {
		log.Printf("Failed fetching opa connection params: %s", err.Error())
		return
	}

	if err = ot.opal.registerCallback(ctx); err != nil {
		log.Printf("Failed registering opal callback: %s", err.Error())
		return
	}

	for {
		// First query without trigger, then wait
		for _, t := range toggles {
			var users []string
			if users, err = opaClient.Query(ctx, t); err != nil {
				// TODO: Retries?
				log.Printf("Failed querying opa for toggle: " + t.Key)
			}
			updates <- ToggleUpdate{Toggle: t, Users: users}
		}

		if err = ot.opal.waitForTrigger(ctx); err != nil {
			return
		}
	}
}

func (ot *OpaTracker) Close() {
	if ot.closeTrack != nil {
		ot.closeTrack()
	}
}

func NewOpaTracker(opalConfig config.OpalConfig) *OpaTracker {
	return &OpaTracker{opal: NewOpalClient(opalConfig)}
}

func TrackAll(ctx context.Context, sources []config.OpalConfig, toggles []config.ToggleConfig) (updates ToggleUpdates) {
	updates = make(ToggleUpdates)

	togglesBySource := make(map[string][]config.ToggleConfig)
	for _, toggle := range toggles {
		srcId := toggle.UsersDocument.Source
		togglesBySource[srcId] = append(togglesBySource[srcId], toggle)
	}

	for _, source := range sources {
		go func(source config.OpalConfig, toggles []config.ToggleConfig, updates ToggleUpdates) {
			tracker := NewOpaTracker(source)

			for {
				// TODO: better retries mechanism
				err := tracker.Track(ctx, toggles, updates)
				log.Printf("Tracking %s ended with error %s", source.Id, err.Error())
			}
		}(source, togglesBySource[source.Id], updates)
	}

	return updates
}

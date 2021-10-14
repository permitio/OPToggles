package trackers

import (
	"context"
	"log"
	"optoggles/config"
)


type OpaTracker struct {
	opal *OpalClient
	closeTrack context.CancelFunc // TODO: is that needed?
}

func (ot *OpaTracker) Track(ctx context.Context, toggles []config.ToggleConfig, results ToggleEvents) (err error) {
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
		if err = ot.opal.waitForTrigger(ctx); err != nil {
			return
		}

		for _, t := range toggles {
			var users []string
			if users, err = opaClient.Query(ctx, t); err != nil {
				// Retries?
				log.Printf("Failed querying opa for toggle: " + t.Name)
			}
			results <- QueryResult{Toggle: t, Users: users}
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

func TrackAll(ctx context.Context, sources []config.OpalConfig, toggles []config.ToggleConfig) (results ToggleEvents) {
	results = make(ToggleEvents)

	togglesBySource := make(map[string][]config.ToggleConfig)
	for _, toggle := range toggles {
		srcId := toggle.UsersDocument.Source
		togglesBySource[srcId] = append(togglesBySource[srcId], toggle)
	}

	for _, source := range sources {
		go func(source config.OpalConfig, toggles []config.ToggleConfig, results ToggleEvents) {
			tracker := NewOpaTracker(source)

			for {
				// TODO: better retries mechanism
				err := tracker.Track(ctx, toggles, results)
				log.Printf("Tracking %s ended with error %s", source.Id, err.Error())
			}
		}(source, togglesBySource[source.Id], results)
	}

	return results
}
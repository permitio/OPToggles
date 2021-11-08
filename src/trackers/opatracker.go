package trackers

import (
	"context"
	"fmt"
	"log"
	"optoggles/config"
	"optoggles/httpserver"
	"optoggles/utils"
)

type OpaTracker struct {
	opal *OpalClient
	opa  *OpaClient
}

func (ot *OpaTracker) Init(ctx context.Context) (err error) {
	if ot.opa, err = ot.opal.getOpaClient(ctx); err != nil {
		log.Printf("Failed fetching opa connection params: %s", err.Error())
		return
	}

	if err = ot.opal.registerCallback(ctx); err != nil {
		log.Printf("Failed registering opal callback: %s", err.Error())
		return
	}

	return nil
}

func (ot *OpaTracker) QueryOnce(ctx context.Context, toggles []config.ToggleConfig, updates ToggleUpdates) error {
	var queryResults []ToggleQueryResult
	for _, t := range toggles {
		if users, err := ot.opa.Query(ctx, t); err != nil {
			return err
		} else {
			queryResults = append(queryResults, ToggleQueryResult{Toggle: t, Users: users})
		}
	}
	updates <- queryResults
	return nil
}

func (ot *OpaTracker) Track(ctx context.Context, toggles []config.ToggleConfig, updates ToggleUpdates) {
	for {
		if err := ot.opal.waitForTrigger(ctx); err != nil {
			// context canceled
			return
		}

		err := utils.Retry(func() error {
			return ot.QueryOnce(ctx, toggles, updates)
		}, fmt.Sprintf("querying source %v", ot.opal.config.Id),
			httpserver.HealthServ.SetPolicyTrackHealth)
		if err != nil {
			return
		}
	}
}

func NewOpaTracker(opalConfig config.OpalConfig) *OpaTracker {
	return &OpaTracker{opal: NewOpalClient(opalConfig)}
}

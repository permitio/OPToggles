package main

import (
	"context"
	"log"
	"optoggles/config"
	"optoggles/httpserver"
	"optoggles/targets"
	"optoggles/trackers"
)

func main() {
	var err error
	var policyTracker *trackers.PolicyTracker
	var targetSync *targets.TargetSync

	ctx := context.Background()

	httpserver.ServeHttp(config.GlobalConfig.Bind)

	if policyTracker, err = trackers.NewPolicyTracker(config.GlobalConfig.Sources, config.GlobalConfig.Toggles); err != nil {
		log.Fatalf("invalid policy source: %v", err)
	}
	if targetSync, err = targets.NewTargetSync(config.GlobalConfig.Target); err != nil {
		log.Fatalf("invalid toggles target: %v", err)
	}

	if err = policyTracker.InitTrackers(ctx); err != nil {
		log.Fatalf("failed to initialize policy tracker: %v", err)
	}
	if err = targetSync.InitTarget(ctx, config.GlobalConfig.Toggles); err != nil {
		log.Fatalf("failed to initialize toggles target: %v", err)
	}

	httpserver.HealthServ.SetStarted()

	policyTracker.TrackAll(ctx)
	err = targetSync.SyncForever(ctx, policyTracker.Updates)
	log.Fatalf("sync forever failed: %v", err)
}

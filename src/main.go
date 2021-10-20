package main

import (
	"context"
	"log"
	"optoggles/config"
	"optoggles/targets"
	"optoggles/trackers"
)

func main() {
	ctx := context.Background()
	target, err := targets.InitTarget(ctx, config.GlobalConfig.Target, config.GlobalConfig.Toggles)
	if err != nil {
		log.Fatalf("failed to initialize toggles target: %sv", err)
	}

	results := trackers.TrackAll(ctx, config.GlobalConfig.Sources, config.GlobalConfig.Toggles)
	log.Printf("Got an error: %b",
		target.Work(ctx, results))
}

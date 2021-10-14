package main

import (
	"context"
	"optoggles/config"
	"optoggles/targets"
	"optoggles/trackers"
)

func main() {
	ctx := context.Background()
	results := trackers.TrackAll(ctx, config.GlobalConfig.Sources, config.GlobalConfig.Toggles)
	targets.NewLogPublisher(results).Work(ctx)
}
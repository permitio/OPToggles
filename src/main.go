package main

import (
	"context"
	"log"
	"optoggles/config"
	"optoggles/ldpublisher"
	"optoggles/opaconnect"
)

func main() {
	ctx := context.Background()
	pub, _ := ldpublisher.NewPrintPublisher()
	queryWorker, err := opaconnect.NewOpaQueryWorker(config.GlobalConfig.OPA.Address,
		config.GlobalConfig.Toggles, pub)

	if err != nil {
		log.Fatalln(err)
	}

	// TODO: Handle connection failures etc...
	go pub.Work(ctx)
	queryWorker.Query(ctx)
	queryWorker.QueryOnTrigger(ctx)
}
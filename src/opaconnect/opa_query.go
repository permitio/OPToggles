package opaconnect

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"optoggles/types"
)

type Publisher interface {
	Publish(usersPerFlag map[string][]string)
}

type OPAQueryWorker struct {
	opaAddress string
	trigger Trigger
	publisher Publisher
	toggles []types.Toggle
}

func (oqw *OPAQueryWorker) Query(ctx context.Context) error {
	log.Printf("OPA querying has been triggered")

	usersPerFlag := make(map[string][]string)
	for _, toggle := range oqw.toggles {
		log.Printf("Querying users in %s for toggle key %s", toggle.GetOpaEndpoint(), toggle.Key)

		// TODO: Use context for post
		// TODO: Use HTTPS
		// TODO: Use authentication
		response, err := http.Get("http://" + oqw.opaAddress + toggle.GetOpaEndpoint())
		if err != nil {
			log.Printf("Can't query OPA: " + err.Error())
			return err
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Can't query OPA: " + err.Error())
			return err
		}

		responseJson := make(map[string][]string)
		if err := json.Unmarshal(responseData, &responseJson); err != nil {
			log.Printf("Can't parse OPA query result: " + err.Error())
			return err
		}

		if users, ok := responseJson["result"]; ok == true {
			usersPerFlag[toggle.Key] = users
		} // TODO: else?
	}

	oqw.publisher.Publish(usersPerFlag)
	return nil
}

func (oqw *OPAQueryWorker) QueryOnTrigger(ctx context.Context) {
	for oqw.trigger.waitForTrigger(ctx) == nil {
		// TODO: ???
		_ = oqw.Query(ctx)
	}
	// Exits when context is canceled
}

func NewOpaQueryWorker(opaAddress string, toggles []types.Toggle, publisher Publisher) (*OPAQueryWorker, error) {
	// TODO: Better design
	cbServer, _ := NewOpalCBServer()
	if err := cbServer.registerCallback(); err != nil {
		return nil, err
	}

	return &OPAQueryWorker{
		opaAddress: opaAddress,
		trigger: cbServer,
		toggles: toggles,
		publisher: publisher,
	}, nil
}
package opaconnect

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Publisher interface {
	Publish(usersPerFlag map[string][]string)
}

type OPAQueryWorker struct {
	opaAddress string
	trigger Trigger
	publisher Publisher
	queries map[string]string
}

func (oqw *OPAQueryWorker) Query(ctx context.Context) error {
	log.Printf("OPA querying has been triggered")

	usersPerFlag := make(map[string][]string)
	for flag, query := range oqw.queries {
		log.Printf("Querying users for toggle " + flag)

		requestBody, err := json.Marshal(map[string]string{"query": query})
		if err != nil {
			log.Printf("Can't marshal query: " + err.Error())
			return err
		}

		// TODO: Use context for post
		// TODO: Use HTTPS
		response, err := http.Post("http://" + oqw.opaAddress + "/v1/query", "application/json",
			bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Can't query OPA: " + err.Error())
			return err
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Can't query OPA: " + err.Error())
			return err
		}

		responseJson := make(map[string][]map[string]string)
		if err := json.Unmarshal(responseData, &responseJson); err != nil {
			log.Printf("Can't parse OPA query result: " + err.Error())
			return err
		}

		if results, ok := responseJson["result"]; ok == true {
			users := make([]string, 0, len(results))
			for _, result := range results {
				// TODO: Understand what is the proper structure
				users = append(users, result["name"])
			}
			usersPerFlag[flag] = users
		}
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

func NewOpaQueryWorker(opaAddress string, queries map[string]string, publisher Publisher) (*OPAQueryWorker, error) {
	// TODO: Better design
	cbServer, _ := NewOpalCBServer()
	if err := cbServer.registerCallback(); err != nil {
		return nil, err
	}

	return &OPAQueryWorker{
		opaAddress: opaAddress,
		trigger: cbServer,
		queries: queries,
		publisher: publisher,
	}, nil
}
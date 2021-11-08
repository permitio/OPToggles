package trackers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"optoggles/config"
	"optoggles/httpserver"
	"optoggles/utils"
)

type Trigger interface {
	waitForTrigger(ctx context.Context) error
}

type OpalClient struct {
	config      config.OpalConfig
	triggerChan chan bool
}

func (oc *OpalClient) getOpaClient(ctx context.Context) (opaClient *OpaClient, err error) {
	log.Printf("Querying opa details from opal")

	var responseData []byte
	if responseData, err = utils.HttpGetBody(ctx, oc.config.Url, "/policy-store/config", oc.config.Token); err != nil {
		return
	}

	responseJson := struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	}{}
	if err = json.Unmarshal(responseData, &responseJson); err != nil {
		log.Printf("Can't parse OPA query result: " + err.Error())
		return
	}
	log.Printf("OPA details. url: %s, token: %s", responseJson.Url, responseJson.Token)

	return NewOpaClient(responseJson.Url, responseJson.Token), nil
}

func (oc *OpalClient) registerCallback(ctx context.Context) (err error) {
	token, err := httpserver.CBServer.AddClient(oc.config.Id, oc.triggerChan)
	if err != nil {
		return err
	}

	// TODO: Use struct for readability?
	callbackRegister := fmt.Sprintf(
		`{"key": "%s", "url": "%s", "config": {"method": "post", "headers": {"Authorization": "Bearer %s"}}}`,
		oc.config.AdvertisedAddress,
		"http://"+oc.config.AdvertisedAddress+httpserver.UpdateEndpoint,
		token)

	return utils.HttpPost(ctx, oc.config.Url, "/callbacks", oc.config.Token, []byte(callbackRegister))
}

func (oc *OpalClient) waitForTrigger(ctx context.Context) error {
	select {
	case <-oc.triggerChan:
		// Empty the trigger channel, avoiding redundant updates
		for i := 0; i < len(oc.triggerChan); i++ {
			_ = <-oc.triggerChan
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func NewOpalClient(config config.OpalConfig) *OpalClient {
	return &OpalClient{
		config:      config,
		triggerChan: make(chan bool, 1024),
	}
}

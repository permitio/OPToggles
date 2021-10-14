package trackers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"optoggles/config"
	"optoggles/utils"
	"strings"
)

type OpaClient struct {
	url string
	token string
}

func getOpaEndpoint(t *config.ToggleConfig) string {
	return "/v1/data/" + strings.ReplaceAll(t.UsersDocument.Package, ".", "/") + "/" + t.UsersDocument.Rule
}

func (oc *OpaClient) Query(ctx context.Context, toggleConfig config.ToggleConfig) (users []string, err error) {
	log.Printf("Querying users in %s for toggle %s", getOpaEndpoint(&toggleConfig), toggleConfig.Name)

	var responseData []byte
	if responseData, err = utils.HttpGetBody(ctx, oc.url, getOpaEndpoint(&toggleConfig), oc.token); err != nil {
		return
	}

	responseJson := make(map[string][]string)
	if err = json.Unmarshal(responseData, &responseJson); err != nil {
		log.Printf("Can't parse OPA query result: " + err.Error())
		return
	}

	var ok bool
	if users, ok = responseJson["result"]; !ok {
		err := errors.New("can't parse OPA query result")
		log.Printf(err.Error())
	}
	return
}

func NewOpaClient(url, token string) *OpaClient {
	return &OpaClient{
		url: url,
		token: token,
	}
}
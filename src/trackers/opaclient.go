package trackers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"optoggles/config"
	"optoggles/utils"
	"strings"
)

type OpaClient struct {
	url   string
	token string
}

func getOpaEndpoint(t *config.ToggleConfig) string {
	return "/v1/data/" + strings.ReplaceAll(t.UsersPolicy.Package, ".", "/") + "/" + t.UsersPolicy.Rule
}

func (oc *OpaClient) Query(ctx context.Context, toggleConfig config.ToggleConfig) (users []string, err error) {
	log.Printf("Querying users in %s for toggle %s", getOpaEndpoint(&toggleConfig), toggleConfig.Key)

	var responseData []byte
	if responseData, err = utils.HttpGetBody(ctx, oc.url, getOpaEndpoint(&toggleConfig), oc.token); err != nil {
		return
	}

	defer func() {
		if err != nil {
			err = errors.New(fmt.Sprintf("malformed OPA query response: %v", err.Error()))
		}
	}()

	responseJson := make(map[string][]string)
	if err = json.Unmarshal(responseData, &responseJson); err != nil {
		return
	}

	var ok bool
	if users, ok = responseJson["result"]; !ok {
		err = errors.New("no result key")
	}
	return
}

func NewOpaClient(url, token string) *OpaClient {
	return &OpaClient{
		url:   url,
		token: token,
	}
}

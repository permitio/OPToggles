package targets

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"optoggles/utils"

	"github.com/mitchellh/mapstructure"
)

type HttpTargetSpec struct {
	EndpointUrl  string
	ExtraHeaders map[string]string
}

type HttpTarget struct {
	HttpTargetSpec
	toggles map[string]map[string]interface{}
}

func NewHttpTarget(spec map[string]interface{}) (*HttpTarget, error) {
	ht := HttpTarget{toggles: make(map[string]map[string]interface{})}
	log.Println(spec)
	if err := mapstructure.Decode(spec, &ht.HttpTargetSpec); err != nil {
		return nil, err
	}
	return &ht, nil
}

func (ht *HttpTarget) doRequest(ctx context.Context, method string, endpoint string, data interface{}) (*http.Response, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(marshaledData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	for headerKey, headerValue := range ht.ExtraHeaders {
		req.Header.Add(headerKey, headerValue)
	}

	return utils.DoRequest(req)
}

func (ht *HttpTarget) CreateToggle(ctx context.Context, key string, spec map[string]interface{}) error {
	if _, keyExists := ht.toggles[key]; keyExists {
		return errors.New("duplicated toggle key")
	}

	toggleData := make(map[string]interface{})
	toggleData["key"] = key
	for k, v := range spec {
		toggleData[k] = v
	}

	res, err := ht.doRequest(ctx, http.MethodPost, ht.EndpointUrl, &toggleData)
	if res != nil && res.StatusCode == http.StatusConflict {
		// If toggle already exists, its content would be patched on the next users update
		log.Printf("toggle %v already exists", key)
	} else if err != nil {
		return err
	}

	ht.toggles[key] = spec
	return nil
}

func (ht *HttpTarget) UpdateToggleWithUsers(ctx context.Context, key string, users []string) error {
	patchData := make(map[string]interface{})
	patchData["users"] = users
	for k, v := range ht.toggles[key] {
		patchData[k] = v
	}

	if _, err := ht.doRequest(ctx, http.MethodPatch, ht.EndpointUrl+"/"+key, &patchData); err != nil {
		return err
	}

	return nil
}

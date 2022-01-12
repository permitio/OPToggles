package targets

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	ldapi "github.com/launchdarkly/api-client-go"
	"github.com/mitchellh/mapstructure"
)

type LaunchdarklyTargetSpec struct {
	LaunchdarklyToken string
}

type LaunchdarklyTarget struct {
	auth    ldapi.APIKey
	client  *ldapi.APIClient
	toggles map[string]LDToggleSpec
}

type LDToggleSpec struct {
	Name         string
	ProjKey      string
	Environments []string
}

func NewLaunchdarklyTarget(spec map[string]interface{}) (*LaunchdarklyTarget, error) {
	ldSpec := LaunchdarklyTargetSpec{}
	if err := mapstructure.Decode(spec, &ldSpec); err != nil {
		return nil, err
	}

	return &LaunchdarklyTarget{
		auth:    ldapi.APIKey{Key: ldSpec.LaunchdarklyToken},
		client:  ldapi.NewAPIClient(ldapi.NewConfiguration()),
		toggles: make(map[string]LDToggleSpec),
	}, nil
}

func (ldt *LaunchdarklyTarget) getContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ldapi.ContextAPIKey, ldt.auth)
}

func (ldt *LaunchdarklyTarget) CreateToggle(ctx context.Context, key string, spec map[string]interface{}) error {
	if _, keyExists := ldt.toggles[key]; keyExists {
		return errors.New("duplicated toggle key")
	}

	var toggleSpec LDToggleSpec
	if err := mapstructure.Decode(spec, &toggleSpec); err != nil {
		return err
	}

	ctx = ldt.getContext(ctx)

	// Flags are created by default as boolean with true/false variations
	_, resp, err := ldt.client.FeatureFlagsApi.PostFeatureFlag(ctx, toggleSpec.ProjKey,
		ldapi.FeatureFlagBody{
			Name:                   toggleSpec.Name,
			Key:                    key,
			ClientSideAvailability: &ldapi.ClientSideAvailability{UsingEnvironmentId: true, UsingMobileKey: true}},
		nil)

	if resp != nil && resp.StatusCode == http.StatusConflict {
		// Flag exists - just update the name
		var name interface{} = toggleSpec.Name
		_, resp, err = ldt.client.FeatureFlagsApi.PatchFeatureFlag(ctx, toggleSpec.ProjKey, key,
			ldapi.PatchComment{Patch: []ldapi.PatchOperation{
				ldapi.PatchOperation{Op: "replace", Path: "/name", Value: &name},
			}})
	}
	if err != nil {
		return err
	}

	ldt.toggles[key] = toggleSpec
	return nil
}

func (ldt *LaunchdarklyTarget) UpdateToggleWithUsers(ctx context.Context, key string, users []string) error {
	toggleSpec := ldt.toggles[key]
	ctx = ldt.getContext(ctx)

	patches := make([]ldapi.PatchOperation, 0)
	for _, env := range toggleSpec.Environments {
		// Turn on targeting per environment,
		var on interface{} = true
		patches = append(patches, ldapi.PatchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/environments/%s/on", env),
			Value: &on,
		})
		// With false as the default variation,
		var ft interface{} = map[string]interface{}{"variation": 1}
		patches = append(patches, ldapi.PatchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/environments/%s/fallthrough", env),
			Value: &ft,
		})
		// And targeting allowed users with the 'true' variation.
		var users interface{} = []map[string]interface{}{{"variation": 0, "values": users}}
		patches = append(patches, ldapi.PatchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/environments/%s/targets", env),
			Value: &users,
		})
	}

	_, _, err := ldt.client.FeatureFlagsApi.PatchFeatureFlag(ctx,
		toggleSpec.ProjKey,
		key,
		ldapi.PatchComment{Patch: patches})

	if err != nil {
		return err
	}
	return nil
}

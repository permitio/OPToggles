package targets

import (
	"context"
	"fmt"
	ldapi "github.com/launchdarkly/api-client-go"
	"log"
	"optoggles/config"
	"optoggles/trackers"
)

type LaunchdarklyTarget struct {
	auth   ldapi.APIKey
	client *ldapi.APIClient
}

func NewLaunchdarklyTarget(token string) *LaunchdarklyTarget {
	return &LaunchdarklyTarget{
		auth:   ldapi.APIKey{Key: token},
		client: ldapi.NewAPIClient(ldapi.NewConfiguration()),
	}
}

func (pp *LaunchdarklyTarget) CreateFlag(ctx context.Context, toggle config.ToggleConfig) error {
	ctx = context.WithValue(ctx, ldapi.ContextAPIKey, pp.auth)

	// Flags are created by default as boolean with true/false variations
	flag, resp, err := pp.client.FeatureFlagsApi.PostFeatureFlag(ctx, toggle.TargetSpec.ProjKey,
		ldapi.FeatureFlagBody{Name: toggle.Name, Key: toggle.TargetSpec.Key}, nil)

	if resp.StatusCode == 409 {
		// Flag exists - just update the name
		var name interface{} = toggle.Name
		flag, resp, err = pp.client.FeatureFlagsApi.PatchFeatureFlag(ctx, toggle.TargetSpec.ProjKey, toggle.TargetSpec.Key,
			ldapi.PatchComment{Patch: []ldapi.PatchOperation{
				ldapi.PatchOperation{Op: "replace", Path: "/name", Value: &name},
			}})
	}
	if err != nil {
		return err
	}

	// Turn on targeting per environment, with false as the default variation
	// TODO: Reuse code
	patches := make([]ldapi.PatchOperation, 0)
	for _, env := range toggle.TargetSpec.Environments {
		var on interface{} = true
		patches = append(patches, ldapi.PatchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/environments/%s/on", env),
			Value: &on,
		})
		var ft interface{} = map[string]interface{}{"variation": 1}
		patches = append(patches, ldapi.PatchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/environments/%s/fallthrough", env),
			Value: &ft,
		})
	}

	flag, resp, err = pp.client.FeatureFlagsApi.PatchFeatureFlag(ctx,
		toggle.TargetSpec.ProjKey,
		toggle.TargetSpec.Key,
		ldapi.PatchComment{Patch: patches})

	if err != nil {
		return err
	}

	fmt.Printf("Created flag: %+v\n", flag)
	return nil
}

func (pp *LaunchdarklyTarget) UpdateFlagTargets(ctx context.Context, result trackers.QueryResult) error {
	ctx = context.WithValue(ctx, ldapi.ContextAPIKey, pp.auth)

	patches := make([]ldapi.PatchOperation, 0)
	for _, env := range result.Toggle.TargetSpec.Environments {
		var users interface{} = []map[string]interface{}{{"variation": 0, "values": result.Users}}
		patches = append(patches, ldapi.PatchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/environments/%s/targets", env),
			Value: &users,
		})
	}

	flag, resp, err := pp.client.FeatureFlagsApi.PatchFeatureFlag(ctx,
		result.Toggle.TargetSpec.ProjKey,
		result.Toggle.TargetSpec.Key,
		ldapi.PatchComment{Patch: patches})
	log.Println(resp, err)

	if err != nil {
		return err
	}
	fmt.Printf("Updated flag: %+v\n", flag)
	return nil
}

func (pp *LaunchdarklyTarget) Work(ctx context.Context, flagsChan trackers.ToggleEvents) error {
	for {
		select {
		case queryResult := <-flagsChan:
			if err := pp.UpdateFlagTargets(ctx, queryResult); err != nil {
				// TODO: ?
				log.Printf("updating flag's target users failed: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

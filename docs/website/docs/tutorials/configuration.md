---
sidebar_position: 4
title: ⚡️ Configuration File
---

### Configuration File structure

`OPToggles` expects its configuration file in `/etc/optoggles/config.yaml`.

The basic structure of the yaml configuration file is:

```yaml
bind: ":8080"

sources:
  - ...
  - ...

target:
  ...

toggles:
  - ...
  - ...
```

| **Path** | **Type** | **Description** |
| :--- | :--- | :--- |
| bind | string | Bind address for OpToggles' HTTP server. This server is use for getting update triggers from `OPAL` and for serving health checks. |

### Sources

The `sources` section is a list of policy sources. Currently only an `OPAL` administrated `OPA` is supported as a policy
source.

Each `source` includes the following attributes

| **Path** | **Type** | **Description** |
| :--- | :--- | :--- |
| id | string | To be referenced from the `toggles` section |
| url | string |  The url base of `Opal Client`. Used to retrieve `OPA`'s connection details and to register the trigger callback|
| token | string | Authorization token for `Opal Client` |
| advertisedAddress | string | Address where this instance of `OPToggles` is reachable to `Opal Client` (for callback registration). |

### Target

The `target` section tells `OPToggles` where it should create and sync its user-authorized toggles to. Supported
targets:

- The [LaunchDarkly](https://launchdarkly.com/) Feature management platform
- Generic REST API HTTP Server - [More Info](https://github.com/permitio/OPToggles/blob/master/example/restapi-config.yaml)

Only one `target` is configured for an instance of `OPToggles`

| **Path** | **Type** | **Description** |
| :--- | :--- | :--- |
| targetType | string | Either `"launchdarkly"` or `"restapi"` |
| targetSpec.launchdarklyToken | string | Required if `targetType` is `"launchdarkly"`. <br/>API access token associated with your LaunchDarkly account. Should have at least `Writer` privileges.  |
| targetSpec.endpointUrl | string | Required if `targetType` is `"restapi"`. <br/>The RestAPI endpoint used to create (/POST) and update (/PATCH) toggles  |
| targetSpec.extraHeaders | map | Optional if `targetType` is `"restapi"`. <br/> Extra headers to include in the REST API requests (e.g. `Authorization`)

### Toggles

The `toggles` section is a list of feature toggles to be managed and continuously updated with the set of allowed
users (queried from a specific rule in one of the policy sources).

Each `toggle` includes the following attributes

| **Path** | **Type** | **Description** |
| :--- | :--- | :--- |
| key | string | Unique identifier for the feature toggle. |
| usersPolicy.source | string | The `id` of the policy source to use  |
| usersPolicy.package | string | The `OPA` package where the desired policy is located  |
| usersPolicy.rule | string | The desired policy's rule name. This rule should return a set of user names for which the feature toggle will be enabled  |
| spec | map | If `targetType` is `"restapi"`, `spec` could contain any user-defined values to be patched to the REST server as part of the toggle object (e.g. `Name`, `Description`, etc...)
| spec.name | string | Required if `targetType` is `"launchdarkly"`. <br/> User readable name for the feature toggle. |
| spec.projKey | string | Required if `targetType` is `"launchdarkly"`. <br/>Key of the LaunchDarkly project under which this toggle should be managed.  |
| spec.environments | list | Required if `targetType` is `"launchdarkly"`. <br/>The environments in which the feature would be enabled for the queried set of users. The feature would be disabled for other environments.   |

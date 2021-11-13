<p  align="center">
 <img src="https://i.ibb.co/BGVBmMK/opal.png" height=170 alt="opal" border="0" />
</p>
<h1 align="center">
OPToggles
</h1>

<h2 align="center">
Open policy based feature toggling
</h2>

<a href="https://hub.docker.com/r/authorizon/optoggles" target="_blank">
    <img src="https://img.shields.io/docker/pulls/authorizon/optoggles?label=Docker%20pulls" alt="Docker pulls">
</a>
<a href="https://opal-access.slack.com/" target="_blank">
    <img src="https://img.shields.io/badge/Slack%20Community-4A154B?logo=slack&logoColor=white" alt="Join our Slack!">
</a>

OPToggles was built to bring the power of OPA into your existing feature toggling solutions. <br/>
It enables you to create user targeted feature flags/toggles based on OPA managed authorization rules.

# Table of contents

- [⚡ Quick Start](#quick-start)
- [🤿 How Does It Work](#design)
- [🦮 First OPToggle Walkthrough](#walkthrough)
- [📖 Configuration Guide](#configuration)
- [👥 Community & Contribution](#community)

# <a name="quick-start"></a>⚡️ Quick Start

OpToggles should be run as a docker container. <br/>
Port 8080 should be exposed to listen to both OPAL trigger callbacks & http health checks. And configuration yaml file
could be supplied through a volume mount.

For example:

```sh
docker run -n optoggles -p 8080:8080 -v $PWD/config.yaml:/optoggles/config.yaml --rm -it authorizon/optoggles:latest
```

Where `config.yaml`:

```yaml
sources:
  - id: myopal
    url: http://opalclient:7000
    advertisedAddress: optoggles:8080

target:
  targetType: launchdarkly
  targetSpec:
    # Replace with your API token
    launchdarklyToken: "api-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

toggles:
  - key: "somefeature"
    usersPolicy:
      source: myopal
      package: "app.rules"
      rule: "somefeature_users"
    spec:
      name: "Some Feature Toggle"
      projKey: "default"
      environments: [ "production", "staging" ]
```

- Optoggles would register to receive callbacks on policy/data changes from `OPAL Client` instance running
  at `opalclient:7000`.
- It would query `OPAL`'s corresponding `OPA` instance for the new value of `somefeature_users` on every
  change (`somefeature_users` is a set of all usernames allowed for some feature).
- It would sync the toggle `somefeature` in your `LaunchDarkly` account to target the current set of usernames allowed
  by the policy. (and create the toggle if it doesn't already exists)
- Health checks are available under `http://optoggles:8080/health[/live,/started]` ([More details](#healthchecks))

# <a name="design"></a>🤿 How Does It Work

### <a name="launchdarkly"></a>OPA & OPAL Integration

### <a name="launchdarkly"></a>The LaunchDarkly Target

### <a name="httpserver"></a>The Generic REST API Target

### <a name="healthchecks"></a>Health Checks

# <a name="walkthrough"></a>🦮 First OPToggle Walkthrough

# <a name="configuration"></a>📖 Configuration Guide

### Basic structure

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
| bind | string | Bind address for OpToggles' HTTP server. This server is use for getting update triggers from `OPAL` & for serving health checks. |

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

The `target` section tells `OPToggles` where it should create & sync its user-authorized toggles to. Supported targets:

- The [LaunchDarkly](https://launchdarkly.com/) Feature management platform
- Generic REST API HTTP Server - [More Info](#httpserver)

Only one `target` is configured for an instance of `OPToggles`

| **Path** | **Type** | **Description** |
| :--- | :--- | :--- |
| targetType | string | Either `"launchdarkly"` or `"http"` |
| targetSpec.launchdarklyToken | string | Required if `targetType` is `"launchdarkly"`. <br/>API access token associated with your LaunchDarkly account. Should have at least `Writer` privileges.  |
| targetSpec.endpointUrl | string | Required if `targetType` is `"http"`. <br/>The RestAPI endpoint used to create (/POST) and update (/PATCH) toggles  |
| targetSpec.extraHeaders | map | Optional if `targetType` is `"http"`. <br/> Extra headers to include in the REST API requests (e.g. `Authorization`)

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
| spec | map | If `targetType` is `"http"`, `spec` could contain any user-defined values to be patched to the REST server as part of the toggle object (e.g. `Name`, `Description`, etc...)
| spec.name | string | Required if `targetType` is `"launchdarkly"`. <br/> User readable name for the feature toggle. |
| spec.projKey | string | Required if `targetType` is `"launchdarkly"`. <br/>Key of the LaunchDarkly project under which this toggle should be managed.  |
| spec.environments | list | Required if `targetType` is `"launchdarkly"`. <br/>The environments in which the feature would be enabled for the queried set of users. The feature would be disabled for other environments.   |

## <a name="community"></a>👥 Community & Contribution

<p  align="center">
 <img src="https://i.ibb.co/m4cN6nW/optoggles.png" height=170 alt="opal" border="0" />
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

`OPToggles` was built to bring the power of `OPA` + `OPAL` into your existing feature toggling solutions. <br/>
It enables you to create user targeted feature flags/toggles based on `OPA` managed authorization rules, and keep them
synced thanks to `OPAL`'s realtime policy and policy-data change detection.<br/>

<img src="docs/optoggles-demo.gif" alt="OPToggles Demo">

# Table of contents

- [⚡ Quick Start](#quick-start)
- [🤿 How Does It Work](#design)
- [🦮 First OPToggle Walkthrough](docs/walkthrough.md)
- [📖 Configuration File Reference](docs/configuration.md)
- [👥 Community & Contribution](#community)

# <a name="quick-start"></a>⚡️ Quick Start

OPToggles should be run as a docker container. <br/>
Port 8080 should be exposed to listen to both OPAL trigger callbacks and http health checks. And configuration yaml file
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

Building your own version of `OPToggles` is as simple as:

```shell
docker build . -t optoggles:$IMAGE_TAG
```

# <a name="design"></a>🤿 How Does It Work

## <a name="opa"></a>OPA + OPAL Integration

OPToggles uses `OPA` as the source of truth regarding what users are authorized for each feature toggle.

`OPAL` is also a necessary component in order to get live policy/data updates and keep your feature toggles always up to
date.

On startup - `OPToggles` would query the configured `OPAL Client` ([See configuration guide](#docs/configuration.md))
for the address of the `OPA` instance it manages.

`OPToggles` would then register to `OPAL Client` as an update callback. `OPAL` would then make an HTTP request
to `OPToggles` on any update to policy or policy data.

The update callback triggers a sync process that would query `OPA` for the current user sets authorized for each feature
toggle. If `OPA` is unavailable, the query would be retried indefinitely with exponential sleep intervals (maxing out at
16 secs).

## <a name="targets"></a>Toggle Targets

The toggle targets are the feature management instances to which `OPToggles` syncs its configured toggles with their
current authorized users lists.

`OPToggles` currerntly supports 2 target types: LaunchDarkly and a Generic REST API target.

[Let us know](#community) if you're interested in other targets.

On startup `OPToggles` creates the toggles, or just updates their metadata if they're already created. Do not manage
your toggles manualy - as your changes would be overridden by `OPToggles`.

On each sync process the target would be updated with the current authorized users for each toggle. On update failure,
it would be retried indefinitely with exponential sleep intervals (maxing out at 16 secs).

For detailed information on configuring each target type see our [configuration guide](#docs/configuration.md).

### <a name="targets"></a>LaunchDarkly

[LaunchDarkly](#https://launchdarkly.com/) is a popular feature management platform, with available SDKs for a wide
range of languages and client side technologies.

Each OpToggle is mapped to a "
Feature Flag" in LaunchDarkly. Feature Flags can have different values across different environments. You can set the
environments you want to be synced in `OPToggles`' configuration.

For each of the set environments of the feature flag, `OPToggles` would turn on "Targeting" mode; with `False` as the
default value, and `True` as the value for each of the authorized (thus targeted) users (queried from `OPA`).

### <a name="targets"></a>Generic REST API

This target type should enable you to integrate with your own custom feature management system.

To create the toggles, `OPToggles` would make a `POST` request to the configured URL. The request's body is a json
representation of the toggle object which contains:

- `key` field with the configured identifier of the toggle
- Other user defined fields (supplied as "key: value" pairs under the toggle's spec)

To update the toggles with the latest authorized users list, `OPToggles` would make a `PATCH` request to "
{configured-url}/{toggle-key}", containing a new value for the `users` field of the toggle.

You can also set user defined HTTP headers that would be added for each REST API call to your target (e.g.
authorization).

## <a name="healthchecks"></a>Health Checks

`OPToggles` exposes a few endpoints to be used for health checking by external services:

- `/health/live` <br/>
  Returns 200 as long as `OPToggles` runs. <br/><br/>
- `/health/started` <br/>
  Returns 200 if init sequence is done (Retrieving OPA's address, registering the update callback & creating the
  toggles). <br/>
  Returns 503 if it isn't done. <br/><br/>
- `/health` <br/>
  Returns 200 if started and running without errors.  
  Returns 500 if there's an ongoing failure to query OPA or update to the target instance.

# <a name="community"></a>👥 Community & Contribution

We would love to hear from you! 😃

**Have a question or a feature request?** <br/>
Post it in our [Github discussions](https://github.com/authorizon/opal/discussions)

**Have an issue with `OPToggles`?** <br/>
Report issues in [Github issues](https://github.com/authorizon/opal/issues)

**Want to chat?** <br/>
Great! join our [Slack community](https://join.slack.com/t/opal-access/shared_invite/zt-nz6yjgnp-RlP9rtOPwO0n0aH_vLbmBQ)

**Want to PM us?** <br/>
Feel free to reach us directy over [email](mailto:or@authorizon.com).

**Interested in contributing?** <br/>
Pull requests are more than welcome!<br/>
Make sure your PR has a corresponding issue.

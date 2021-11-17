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

`OPToggles` was built to bring the power of `OPA` + `OPAL` into your existing feature toggling solutions. <br/>
It enables you to create user targeted feature flags/toggles based on `OPA` managed authorization rules, and keep them
synced thanks to `OPAL`'s realtime policy and policy-data change detection.<br/>

# Table of contents

- [⚡ Quick Start](#quick-start)
- [🤿 How Does It Work](#design)
- [🦮 First OPToggle Walkthrough](docs/walkthrough.md)
- [📖 Configuration File Reference](docs/configuration.md)
- [👥 Community & Contribution](#community)

# <a name="quick-start"></a>⚡️ Quick Start

OpToggles should be run as a docker container. <br/>
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

# <a name="design"></a>🤿 How Does It Work

### <a name="launchdarkly"></a>OPA + OPAL Integration

### <a name="launchdarkly"></a>The LaunchDarkly Target

### <a name="httpserver"></a>The Generic REST API Target

### <a name="healthchecks"></a>Health Checks

## <a name="community"></a>👥 Community & Contribution

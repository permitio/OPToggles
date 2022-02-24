# <a name="intro"></a>⚡️ Quick Start Guide to OPToggles
## 

OPToggles should be run as a docker container. <br/>

Port 8080 should be exposed to listen to both OPAL trigger callbacks and http health checks. And configuration yaml file
could be supplied through a volume mount.

For example:

```sh
docker run -n optoggles -p 8080:8080 -v $PWD/config.yaml:/optoggles/config.yaml --rm -it permitio/optoggles:latest
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
- Health checks are available under `http://optoggles:8080/health[/live,/started]` ([More details](howitworks.md#healthchecks))

Building your own version of `OPToggles` is as simple as:

```shell
docker build . -t optoggles:$IMAGE_TAG
```
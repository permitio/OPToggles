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

`OPToggles` was built to bring the power of `OPA` & `OPAL` into your existing feature toggling solutions. <br/>
It enables you to create user targeted feature flags/toggles based on `OPA` managed authorization rules, and keep them
synced thanks to `OPAL`'s realtime policy and policy-data change detection.<br/>

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

Let's walk you through setting up your first open policy based feature toggles, integrated with `LaunchDarkly`!

### Setup `OPA` & `OPAL`

We'll use `docker-compose` for the task, the end result is
available [here](https://github.com/authorizon/OPToggles/tree/master/example).

First things first, we're going to need an `OPA` instance, managed by `OPAL` for realtime policy and policy data
updates. Our starting point would be one of `OPAL`'s
example [docker-compose configurations](https://github.com/authorizon/opal/blob/master/docker/docker-compose-example.yml)
.<br/>(To learn more about working with `OPAL` container images
view [this guide](https://github.com/authorizon/opal/blob/master/docs/HOWTO/get_started_with_opal_using_docker.md))

In order to have policy based user-targeted feature toggles, we're going to need some users data, and some policies:

We'll get `data.json` from this `Authorizon`'
s [example policy repo](https://github.com/authorizon/opal-example-policy-repo/blob/master/data.json). And we'll create
a new file `features.rego` with the following rego rules:

```rego
package app.rbac

billing_users[users]{
  some user, i
  data.users[user].roles[i] == "billing"
  users := user
}

us_users[users]{
  some user
  data.users[user].location.country == "US"
  users := user
}
```

This snippet declares two sets of users: `billing_users` returns a set of all users that have the billing
role. `us_users` returns a set of all users which are located in the US.

We want to feed `OPA` with this data & policies, for that - we'll use `OPAL`'
s [git-tracking capabilities](https://github.com/authorizon/opal/blob/master/docs/HOWTO/track_a_git_repo.md):

1. Put both the policy & data files in a git repo with a `.manifest`  file listing the files paths `OPAL` need to track.
   We already have our example files in this repo
   under [example/policy](https://github.com/authorizon/OPToggles/tree/master/example/policy).
2. Edit our `docker-compose.yaml` to configure `OPAL Server` to track the right git repo, branch & `.manifest` file:

```yaml
  opal_server:
    environment:
      - OPAL_POLICY_REPO_URL=https://github.com/authorizon/OPToggles
      - OPAL_POLICY_REPO_MAIN_BRANCH=master
      - OPAL_POLICY_REPO_MANIFEST_PATH=.manifest
      - OPAL_POLICY_REPO_POLLING_INTERVAL=30
```

So we've got policies in place defining the set of usernames allowed using certain features. You can already use `OPAL`
with your backend to authorize requests using realtime data, and deny users forbidden actions.

But getting an `401 Unauthorized` (or another error message, as elegant as it might be) in the client side isn't exactly
an UX best practice :)<br/>

### Setup `LaunchDarkly`

That's where the magic of feature management platforms enters: `LaunchDarkly` enables you to manage feature toggles
across multiple projects and deployment environments, and it has rich client-side sdk support, already used by many
developers to turn UI features on & off.

If you don't already have a `LaunchDarkly` account, create it [here](https://app.launchdarkly.com/signup).

Next thing you would need is a `project` and one or more `environment`. You can manage those under `Account settings`
-> `Projects`.

Each `LaunchDarkly` account can contain multiple projects for different products. I this example we're going to use the
pre-existing `default` project. If you want to create a new one - do
it [here](https://app.launchdarkly.com/settings/projects/new).

Your feature toggles can have different configurations for different deployment environments, you would need at least
one environment, which you can create [here](https://app.launchdarkly.com/settings/projects/default/env/new). In our
example - we defined 2 environments: `production` & `test`.

Everything should look like that:

<img src="https://i.ibb.co/8sZ43bp/image.jpg" alt="LaunchDarkly Project Settings"/>

From now on - OPToggles would take care of toggle creation and updates thorugh `LaunchDarkly` for us. In order to so it
needs an access token. You can create a token under `Account settings` -> `Authorization`, or
using [this link](https://app.launchdarkly.com/settings/authorization/tokens/new). For `OPToggles` to work, the token
should have `Writer` permissions. <br/>
Don't loose the generated token! We're gonna need it soon.

### Setup `OPToggles`

Now let's bring everything together using `OPToggles`. Our configuration yaml
is [here](https://github.com/authorizon/OPToggles/blob/master/example/launchdarkly-config.yaml), and looks like that:

```yaml
sources:
  - id: example-opal
    url: http://opal_client:7000
    token: ""
    advertisedAddress: optoggles:8080

target:
  targetType: launchdarkly
  targetSpec:
    # Replace with your generated api token
    launchdarklyToken: "api-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

toggles:
  - key: "us-feature"
    usersPolicy:
      source: example-opal
      package: "app.rbac"
      rule: "us_users"
    spec:
      name: "US Only Feature"
      projKey: "default"
      environments: [ "production", "test" ]
  - key: "billing-feature"
    usersPolicy:
      source: example-opal
      package: "app.rbac"
      rule: "billing_users"
    spec:
      name: "Billing Feature"
      projKey: "default"
      environments: [ "test" ]
```

Important comments:

1. Since we're running everything in docker-compose, we simply use the service names as hostnames (`opal_client`
   , `optoggles`)
2. If your `OPAL Client` runs in secure mode, please supply an JWT authentication token.
3. Each toggle's `usersPolicy` defines the `OPA` source - the package & rule strings match the names from
   the `features.rego` file.
4. `OPToggles` would query the `OPA` instance that is associated with the supplied `OPAL CLient`.
5. View our [configuration guide](#configuration) for full understanding of its format.

Now we can add the `OPToggles` as a service to
our [docker-compose](https://github.com/authorizon/OPToggles/blob/master/example/docker-compose.yaml):

```yaml
  optoggles:
    image: authorizon/optoggles:latest
    depends_on:
      - opal_client
    restart: on-failure
    volumes:
      - $PWD/launchdarkly-config.yaml:/etc/optoggles/config.yaml
```

Setting `restart` to `on-failure` is comfortable solution to handle errors on `OPToggles` initiation when rest of the
containers are not fully started.

Now that we have everything in place - let's bring it up!

```shell
docker-compose up -d
```

### Consuming the feature toggles

If everything went well. you should see the newly created flags in your `LaunchDarkly` account:

Let's update bob's location to another country:

```shell
opal-client publish-data-update --src-url https://api.country.is/23.54.6.78 -t policy_data --dst-path /users/bob/locationgit:master*
```

And our "US Only Feature" should be immediately updated to exclude "bob"!

You can integrate these new toggles into your client-side code like you would with any other `LaunchDarkly` flag. If
that's your first time - https://docs.launchdarkly.com/sdk/client-side.

One last important note: the created toggles are fully managed by `OPToggles`, trying to make manual changes to them
makes no sense as they would get overridden by `OPToggles`.

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

# <a name="walkthrough"></a>🦮 First OPToggle Walkthrough

Let's walk you through setting up your first open policy based feature toggles, integrated with `LaunchDarkly`!

<br/>

## Setup `OPA` + `OPAL`

We'll use `docker-compose` for the task, the end result is
available [here](https://github.com/permitio/OPToggles/tree/master/example).

First things first, we're going to need an `OPA` instance, managed by `OPAL` for realtime policy and policy data
updates. Our starting point would be one of `OPAL`'s
example [docker-compose configurations](https://github.com/permitio/opal/blob/master/docker/docker-compose-example.yml)
.<br/>(To learn more about working with `OPAL` container images
view [this guide](https://github.com/permitio/opal/blob/master/docs/HOWTO/get_started_with_opal_using_docker.md))

In order to have policy based user-targeted feature toggles, we're going to need some users data, and some policies:

We'll get `data.json` from this `Permit.io`'
s [example policy repo](https://github.com/permitio/opal-example-policy-repo/blob/master/data.json). And we'll create a
new file `features.rego` with the following rego rules:

```rego
package app.rbac

billing_users[users]{
  some user, i
  data.example.users[user].roles[i] == "billing"
  users := user
}

us_users[users]{
  some user
  data.example.users[user].location.country == "US"
  users := user
}
```

This snippet declares two sets of users: `billing_users` returns a set of all users that have the billing
role. `us_users` returns a set of all users which are located in the US.

We want to feed `OPA` with that data & policies, for that - we'll use `OPAL`'
s [git-tracking capabilities](https://github.com/permitio/opal/blob/master/docs/HOWTO/track_a_git_repo.md):

1. Put both the policy & data files in a git repo with a `.manifest`  file listing the files paths `OPAL` need to track.
   We already have our example files in this repo
   under [example/policy](https://github.com/permitio/OPToggles/tree/master/example/policy).
2. Edit our `docker-compose.yaml` to configure `OPAL Server` to track the right git repo, branch & `.manifest` file:

```yaml
  opal_server:
    environment:
      - OPAL_POLICY_REPO_URL=https://github.com/permitio/OPToggles
      - OPAL_POLICY_REPO_MAIN_BRANCH=master
      - OPAL_POLICY_REPO_MANIFEST_PATH=example/.manifest
      - OPAL_POLICY_REPO_POLLING_INTERVAL=30
```

So we've got policies in place defining the set of usernames allowed using certain features. You can already use `OPAL`
with your backend to authorize requests using realtime data, and deny users forbidden actions.

But getting an `401 Unauthorized` (or another error message, as elegant as it might be) in the client side isn't exactly
an UX best practice :)

<br/>

## Setup `LaunchDarkly`

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
one environment, which you can create [here](https://app.launchdarkly.com/settings/projects/default/env/new). <br/>
In our example - we use the 2 pre-defined environments: `production` & `test`.

Everything should look like that:

<img src="https://i.ibb.co/8sZ43bp/image.jpg" alt="LaunchDarkly Project Settings"/>

From now on - OPToggles would take care of toggle creation and updates thorugh `LaunchDarkly` for us. <br/>
**In order to so it needs an access token (aka API key), those are only available for Professional plans or higher at
the moment. <br/>**
You can create a token under `Account settings` -> `Authorization`, or
using [this link](https://app.launchdarkly.com/settings/authorization/tokens/new). <br/>
`OPToggles` requires your token to have `Writer` permissions at minimum. <br/>
Don't loose the generated token! We're gonna need it soon.

### Setup `OPToggles`

Now let's bring everything together using `OPToggles`. Our configuration yaml
is [here](https://github.com/permitio/OPToggles/blob/master/example/launchdarkly-config.yaml), and looks like that:

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
5. View our [configuration guide](configuration.md) for full understanding of its format.

Now we can add the `OPToggles` as a service to
our [docker-compose](https://github.com/permitio/OPToggles/blob/master/example/docker-compose.yaml):

```yaml
  optoggles:
    image: permitio/optoggles:latest
    depends_on:
      - opal_client
    restart: on-failure
    volumes:
      - $PWD/launchdarkly-config.yaml:/etc/optoggles/config.yaml
```

Setting `restart` to `on-failure` is useful for errors on `OPToggles` initiation when rest of the containers are not
fully started.

Now that we have everything in place - let's bring it up!

```shell
docker-compose up -d
```

### Consuming the feature toggles

If everything went well. you should see the newly created flags in your `LaunchDarkly` account:

<img src="https://i.ibb.co/QncsFG7/toggle-before.png" alt="toggle-before" border="0">

Let's update bob's location to another country:

```shell
opal-client publish-data-update --src-url https://api.country.is/23.54.6.78 -t policy_data --dst-path /users/bob/locationgit:master*
```

And our "US Only Feature" should be immediately updated to exclude "bob"!

<img src="https://i.ibb.co/d2P71BJ/toggle-after.png" alt="toggle-after"/>


You can integrate these new toggles into your client-side code like you would with any other `LaunchDarkly` flag. If
that's your first time - https://docs.launchdarkly.com/sdk/client-side.

One last important note: the created toggles are fully managed by `OPToggles`, trying to make manual changes to them
makes no sense as they would get overridden by `OPToggles`.

## Description

This plugin enables the ability to manage artifacts in [npm](https://www.npmjs.org/) in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-npm

Registry: https://hub.docker.com/r/target/vela-npm

## Usage

> **NOTE:**
>
> Users should refrain from using latest as the tag for the Docker image.
>
> It is recommended to use a semantically versioned tag instead.

Sample of publishing package:

```yaml
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    secrets: [ npm_password ]
    parameters:
      username: npmUsername
      registry: https://registry.npmjs.org
```

Sample of publishing if registry does not support `npm ping`:

{{% notice tip %}}
Recommended if you are deploying to a registry inside ***Artifactory***
{{% /notice %}}

```diff
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    secrets: [ npm_password ]
    parameters:
      username: npmUsername
      registry: https://registry.npmjs.org
+     skip_ping: true
```

Sample of pretending to publish package:

```diff
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    secrets: [ npm_password ]
    parameters:
      username: npmUsername
      registry: https://registry.npmjs.org
+     dry_run: true
```

Sample of first time publishing package:

```diff
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    secrets: [ npm_password ]
    parameters:
      username: npmUsername
      registry: https://registry.npmjs.org
+     first_publish: true
```

Sample of publishing with additional dist-tag:

{{% notice warning %}}
Tags are used as an alias and cannot be valid semver
{{% /notice %}}

```diff
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    secrets: [ npm_password ]
    parameters:
      username: npmUsername
      registry: https://registry.npmjs.org
+     tag: beta
```

Higher level of tolerance for npm audit:

```diff
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    secrets: [ npm_password ]
    parameters:
      username: npmUsername
      registry: https://registry.npmjs.org
+     audit_level: critical
```

## Secrets

{{% notice warning %}}
**Users should refrain from configuring sensitive information in their pipeline in plain text.**
{{% /notice %}}

### Internal

The plugin accepts the following `parameters` for authentication:

| Parameter  | Environment Variable Configuration    |
| ---------- | ------------------------------------- |
| `password` | `NPM_PASSWORD`, `PARAMETER_PASSWORD`  |
| `username` | `NPM_USERNAME`, `PARAMETER_USERNAME`  |
| `registry` | `NPM_REGISTRY`, `PARAMETER_REGISTRY`  |
| `email`    | `NPM_EMAIL`, `PARAMETER_EMAIL`        |

Users can use [Vela internal secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    secrets: [ npm_password ]
    parameters:
      username: npmUsername
      registry: https://registry.npmjs.org
-     password: superSecretPassword
```


> This example will add the `secrets` to the `npm_publish` step as environment variables:
> - `NPM_PASSWORD`=value

### External


The plugin accepts the following files for authentication:

| Parameter  | Volume Configuration                                          |
| ---------- | ------------------------------------------------------------- |
| `password` | `/vela/parameters/npm/password`, `/vela/secrets/npm/password` |
| `username` | `/vela/parameters/npm/username`, `/vela/secrets/npm/username` |
| `registry` | `/vela/parameters/npm/registry`, `/vela/secrets/npm/registry` |
| `email`    | `/vela/parameters/npm/email`, `/vela/secrets/npm/email`       |

Users can use [Vela external secrets](https://go-vela.github.io/docs/concepts/pipeline/secrets/origin/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: npm_publish
    image: target/vela-npm:latest
    pull: not_present
    parameters:
      registry: https://registry.npmjs.org
-     username: npmUsername
-     password: superSecretPassword
```

## Parameters

The following parameters are used to configure the image:

| Name            | Description                                                                                                        | Required | Default                      |
| --------------- | ------------------------------------------------------------------------------------------------------------------ | -------- | ---------------------------- |
| `username`      | username for communication with npm                                                                                | `true`   | `N/A`                        |
| `password`      | password for communication with npm                                                                                | `false`  | `N/A`                        |
| `email`         | email for communication with npm                                                                                   | `false`  | `N/A`                        |
| `registry`      | npm instance to communicate with                                                                                   | `false`  | `https://registry.npmjs.org` |
| `audit_level`   | level at which the audit check should fail (valid options: `low`, `moderate`, `high`, `critical`, `none` to skip)  | `false`  | `low`                        |
| `strict_ssl`    | whether or not to do SSL key validation during communication                                                       | `false`  | `true`                       |
| `always_auth`   | force npm to always require authentication                                                                         | `false`  | `false`                      |
| `skip_ping`     | whether or not to skip `npm ping` authentication command                                                           | `false`  | `false`                      |
| `dry_run`       | enables pretending to perform the action                                                                           | `false`  | `false`                      |
| `tag`           | publish package with given alias tag                                                                               | `false`  | `latest`                     |
| `log_level`     | set the log level for the plugin (valid options: `info`, `debug`, `trace`)                                         | `true`   | `info`                       |

## package.json
This is your module's manifest.  There are a few important keys that need to be set in order to publish your module

* **name** - your package name that will be checked against in the registry
* **version** - your package version that will be used to publish, it must be valid semver and unique to the registry
* **private** - this needs to be set to `false` even if you are publishing it internally.
* **publishConfig** - this should be configured to your registry location and registry parameter should match this value

For example values, see npm's [documentation](https://docs.npmjs.com/files/package.json)

## Template

COMING SOON!

## Troubleshooting

{{% notice tip %}}
**Here are the available log levels to assist in troubleshooting:**
trace, debug, info, warn, error, fatal, panic
{{% /notice %}}
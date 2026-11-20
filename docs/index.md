---
layout: "zscaler"
page_title: "Provider: Zero Trust Cloud (ZTC)"
description: |-
    The Zero Trust Cloud  provider is used to interact with Zero Trust Cloud (ZTC) API
---

# Zero Trust Cloud (ZTC) Provider

The Zero Trust Cloud  provider is used to interact with ZTC API, to automate the provisioning of new locations, IPSec and GRE tunnels, URL filtering policies, Cloud Firewall Policies, DLP Dictionaries, Local Accounts etc. The provider is intended to save time and reducing configuration errors. With this ZTC provider, DevOps teams can automate their security and transform it into DevSecOps workflows. To use this  provider, you must create ZTC API credentials.

Use the navigation on the left to read about the available resources.

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](guides/troubleshooting.md) for guidance on typical problems.

## Feature Availability and API Parity

-> **Important:** The ZTC Terraform provider maintain parity with publicly available API endpoints. In some instances, certain features or attributes available via the Zscaler UI may not be immediately available through the API, and therefore cannot be included in the Terraform provider. This does not indicate that the provider is lagging behind; rather, it reflects that we implement only the features that are currently exposed by the public API.

If there is a feature or attribute you would like to see included in the provider, you are welcome to:

- Submit a feature request via [GitHub Issues](https://github.com/zscaler/terraform-provider-ztc/issues)
- Contact Zscaler Global Support by opening a support ticket.

Our team continuously works with product teams to expand API coverage and will incorporate new features into the provider as they become publicly available through the API.

## Zscaler OneAPI New Framework

The ZTC Terraform Provider now offers support for [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) Oauth2 authentication through [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE** As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE**: Attention Government customers. OneAPI and Zidentity support the government (FedRAMP) clouds via the unified `cloud=gov` and `cloud=govus` values, but this capability is **not yet available in a published release of this provider**. Authenticate government tenants using the legacy API framework for now. See [OneAPI Government (FedRAMP) Cloud Environments](#oneapi-government-fedramp-cloud-environments) below.

## Examples Usage - Client Secret Authentication

```hcl
# Configure the Zero Trust Cloud Provider
terraform {
    required_providers {
        ztc = {
            version = "~> 0.1.6"
            source = "zscaler/ztc"
        }
    }
}

# Configure the ZTC Provider (OneAPI Authentication)
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "ztc" {
  client_id = "[ZSCALER_CLIENT_ID]"
  client_secret = "[ZSCALER_CLIENT_SECRET]"
  vanity_domain = "[ZSCALER_VANITY_DOMAIN]"
  zscaler_cloud = "[ZSCALER_CLOUD]"
}
```

## Examples Usage - Private Key Authentication

```hcl
# Configure the Zero Trust Cloud Provider
terraform {
    required_providers {
        ztc = {
            version = "~> 0.1.6"
            source = "zscaler/ztc"
        }
    }
}

# Configure the ZTC Provider (OneAPI Authentication) - Private Key
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "ztc" {
  client_id     = "[ZSCALER_CLIENT_ID]"
  private_key   = "[ZSCALER_PRIVATE_KEY]"
  vanity_domain = "[ZSCALER_VANITY_DOMAIN]"
  zscaler_cloud = "[ZSCALER_CLOUD]"
}
```

**NOTE**: The `zscaler_cloud` is optional and only required when authenticating to other environments i.e `beta`. Notice that ONLY the value `beta` is supported.

**NOTE**: The `zscaler_cloud` is optional for production comercial Clouds and ONLY required when authenticating to other environments. Currently the only supported value are:

- Production Commericial Clouds: The cloud parameter `IS NOT` required.
- Test (Beta) Commercial Clouds: `beta`
- FedRAMP Clouds: `gov` or `govus` — not yet available in a published release; see the note below

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

For the resources and data sources examples, please check the [examples](https://github.com/zscaler/terraform-provider-ztc/tree/master/examples) directory.

## Authentication - OneAPI New Framework

This provider supports authentication via the new Zscaler API framework [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi)

Zscaler OneAPI uses the OAuth 2.0 authorization framework to provide secure access to Zero Trust Cloud (ZTC) APIs. OAuth 2.0 allows third-party applications to obtain controlled access to protected resources using access tokens. OneAPI uses the Client Credentials OAuth flow, in which client applications can exchange their credentials with the authorization server for an access token and obtain access to the API resources, without any user authentication involved in the process.

* [ZTC API](https://help.zscaler.com/cloud-branch-connector/getting-started-cloud-branch-connector-api)

### Default Environment variables

You can provide credentials via the `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET`, `ZSCALER_VANITY_DOMAIN`, `ZSCALER_CLOUD` environment variables, representing your Zidentity OneAPI credentials `clientId`, `clientSecret`, `vanityDomain` and `zscaler_cloud` respectively.

| Argument        | Description                                                                                         | Environment Variable     |
|-----------------|-----------------------------------------------------------------------------------------------------|--------------------------|
| `client_id`     | _(String)_ Zscaler API Client ID, used with `clientSecret` or `PrivateKey` OAuth auth mode.         | `ZSCALER_CLIENT_ID`      |
| `client_secret` | _(String)_ Secret key associated with the API Client ID for authentication.                         | `ZSCALER_CLIENT_SECRET`  |
| `privateKey`    | _(String)_ A string Private key value.                                                              | `ZSCALER_PRIVATE_KEY`    |
| `vanity_domain` | _(String)_ Refers to the domain name used by your organization.                                     | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud`         | _(String)_ The name of the Zidentity cloud, e.g., beta.                                             | `ZSCALER_CLOUD`          |

### Alternative OneAPI Cloud Environments

OneAPI supports authentication and can interact with alternative Zscaler enviornments i.e `beta`. To authenticate to these environments you must provide the following values:

| Argument         | Description                                                                                         |   | Environment Variable     |
|------------------|-----------------------------------------------------------------------------------------------------|---|--------------------------|
| `vanity_domain`   | _(String)_ Refers to the domain name used by your organization |   | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud`          | _(String)_ The name of the Zidentity cloud i.e beta      |   | `ZSCALER_CLOUD`          |

For example: Authenticating to Zscaler Beta environment:

```sh
export ZSCALER_VANITY_DOMAIN="acme"
export ZSCALER_CLOUD="beta"
```

## OneAPI Government (FedRAMP) Cloud Environments

OneAPI supports the Zscaler government (FedRAMP) clouds. These are FedRAMP-isolated environments served by a dedicated Zidentity identity provider and API gateway. To authenticate, set the `cloud` attribute (or `ZSCALER_CLOUD` environment variable) to one of the supported government values:

| Argument        | Description                                                         | Environment Variable     |
|-----------------|---------------------------------------------------------------------|--------------------------|
| `vanity_domain` | _(String)_ Refers to the domain name used by your organization     | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud` | _(String)_ Supported Zidentity Gov Cloud `gov` or `govus`                | `ZSCALER_CLOUD`          |

~> **NOTE** FedRAMP cloud support is not yet available in a published release of this provider. The latest release, `v0.2.1`, does not include it. Until a release that does is available, authenticate government tenants using the legacy client with `ztc_cloud = "zscalergov"` or `"zscalerten"`. See [Authentication - Legacy API Framework](#authentication---legacy-api-framework).

For example, authenticating to the GOV environment:

```sh
export ZSCALER_VANITY_DOMAIN="acme"
export ZSCALER_CLOUD="gov"
```

### OneAPI (API Client Scope)

OneAPI Resources are automatically created within the ZIdentity Admin UI based on the RBAC Roles
applicable to APIs within the various products. For example, in ZTC, navigate to `Administration -> Role
Management` and select `Add API Role`.

Once this role has been saved, return to the ZIdentity Admin UI and from the Integration menu
select API Resources. Click the `View` icon to the right of Zscaler APIs and under the ZTC
dropdown you will see the newly created Role. In the event a newly created role is not seen in the
ZIdentity Admin UI a `Sync Now` button is provided in the API Resources menu which will initiate an
on-demand sync of newly created roles.

## Legacy API Framework

### ZTC native authentication

* As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

### Examples Usage

```hcl
# Configure the Zero Trust Cloud  Provider
terraform {
    required_providers {
        ztc = {
            version = "~> 0.1.0"
            source = "zscaler/ztc"
        }
    }
}

# Configure the ZTC Provider (Legacy Authentication)
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "ztc" {
  username            = "[ZTC_USERNAME]"
  password            = "[ZTC_PASSWORD]"
  api_key             = "[ZTC_API_KEY]"
  ztc_cloud           = "[ZTC_CLOUD]"
  use_legacy_client   = "[ZSCALER_USE_LEGACY_CLIENT]"
}
```

The ZTC Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `zscaler`
* `zscloud`
* `zscalerone`
* `zscalertwo`
* `zscalerthree`
* `zscalerbeta`
* `zscalergov`
* `zscalerten`
* `zspreview`

### Environment variables

You can provide credentials via the `ZTC_USERNAME`, `ZTC_PASSWORD`, `ZTC_API_KEY`, `ZTC_CLOUD`, `ZSCALER_USE_LEGACY_CLIENT` environment variables, representing your ZTC `username`, `password`, `api_key`,  `ztc_cloud` and `use_legacy_client` respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `username`       | _(String)_ A string that contains the email ID of the API admin.| `ZTC_USERNAME` |
| `password`       | _(String)_ A string that contains the password for the API admin.| `ZTC_PASSWORD` |
| `api_key`       | _(String)_ A string that contains the obfuscated API key (i.e., the return value of the obfuscateApiKey() method).| `ZTC_API_KEY` |
| `ztc_cloud`       | _(String)_ The host and basePath for the cloud services API is `$zsapi.<Zscaler Cloud Name>/api/v1`.| `ZTC_CLOUD` |
| `use_legacy_client`       | _(Bool)_ Enable use of the legacy ZTC API Client.| `ZSCALER_USE_LEGACY_CLIENT` |

```sh
# Change place holder values denoted by brackets to real values, including the
# brackets.

$ export ZTC_USERNAME="[ZTC_USERNAME]"
$ export ZTC_PASSWORD="[ZTC_PASSWORD]"
$ export ZTC_API_KEY="[ZTC_API_KEY]"
$ export ZTC_CLOUD="[ZTC_CLOUD]"
$ export ZSCALER_USE_LEGACY_CLIENT=true
$ terraform plan
```

If you are on Windows, use PowerShell to set the environmenr variables using the following commands:

```pwsh
$env:username = 'xxxxxxxxxxxxxxxx'
$env:password = 'xxxxxxxxxxxxxxxx'
$env:api_key = 'xxxxxxxxxxxxxxxx'
$env:ztc_cloud = '<zscaler_cloud_name>'
$env:use_legacy_client = true
```

```hcl
# provider settings established with values from environment variables
provider "ztc" {}
```

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

For details about how to retrieve your tenant Base URL and API key/token refer to the Zscaler help portal. <https://help.zscaler.com/cloud-branch-connector/getting-started-cloud-branch-connector-api>

## ZTC Configuration Activation

The ZTC platform saves configuration changes in a **pending** state. A successful `terraform apply` writes the configuration to the tenant but does not put it into effect — the changes must be **activated**.

Activation is a tenant-wide operation, not a per-resource one. **A single activation call publishes every pending change in the tenant**, regardless of how many resources Terraform created, updated, or deleted.

### Method 1 — Out-of-band activation with the `ztcActivator` CLI (recommended)

Compile the dedicated activation binary and run it once, after Terraform has finished:

```bash
terraform apply && ztcActivator
```

This is the recommended method for CI/CD pipelines and for large configurations. It issues exactly one activation call for the entire run and keeps activation timing explicitly under your control. See [ZTC Activator](guides/ztc-activator-overview.md) for build and usage details.

### Method 2 — The `ztc_activation_status` resource

Declare the resource and use the `depends_on` meta-argument so that it runs after the resources it should activate. Because `depends_on` cannot be inferred, **every resource whose changes must be activated has to be listed explicitly** — anything you omit may be applied after activation and left pending. In configurations built from reusable modules that list becomes difficult to maintain, so prefer Method 1 there.

### API Session Timeout and Long-Running Applies

As with ZIA, the platform activates pending changes when a session ends, including when an API session reaches its timeout. Refer to the [Zscaler API documentation](https://help.zscaler.com/legacy-apis/configuring-postman-rest-api-client-3) for the session behaviour of the legacy API framework.

!> **IMPORTANT:** Unlike ZIA — which exposes an adjustable API session timeout of 5 to 20 minutes — **ZTC provides no way to adjust the API session lifetime**, either through this provider or through the tenant's settings. The session length is fixed by the platform and cannot be raised.

This makes run duration the only lever available. A run that outlives its session will have pending changes activated part-way through, before the apply has finished writing the rest of the configuration; the provider then re-authenticates and continues, so the run does not fail, but activation has already happened at a point you did not choose. Keep individual runs short by splitting large configurations across smaller states, and prefer the out-of-band activator so that the activation you *intend* is an explicit, separate step.

### Running Several Configurations Against One Tenant

Splitting a large deployment across several Terraform configurations is fully supported, and the previous section recommends it as the way to keep runs short. Two operational rules apply, because a ZTC tenant has a single configuration-change lock and a single activation queue regardless of how state is split:

- **Apply one configuration at a time** against a given tenant. `terraform plan` can run in parallel — reads are unaffected. In CI, key the concurrency control on the tenant. Concurrent applies surface as `EDIT_LOCK_NOT_AVAILABLE` or `Failed during enter Org barrier`; the provider retries these, so the usual symptom is a slow run rather than an error.
- **Activate once, after the last apply** — not once per configuration. Activation publishes every pending change in the tenant, so a per-configuration call is redundant and will queue behind any other session that still has unactivated changes.

~> **NOTE:** Serialising applies lengthens total pipeline time, so keep each configuration small enough that its apply finishes comfortably inside an API session. Because ZTC's session lifetime cannot be raised, this matters more here than on ZIA.

## Rate Limiting

The ZTC platform enforces API rate limits on a per-endpoint basis, and different endpoints have different thresholds.

**You do not need to configure anything to stay within these limits.** When a limit is exceeded, the API returns an HTTP `429 Too Many Requests` response with a `Retry-After` header, and the provider waits the indicated interval and retries the request transparently. Exceeding a rate limit therefore costs a short delay rather than a failed apply, including on large bulk deployments.

Run Terraform with its default settings. In particular, do not reduce Terraform's `-parallelism`: the flag applies to the entire run and cannot be scoped to a single resource type or endpoint, so lowering it slows every resource in the configuration. Rule-based resources are affected most, because rule placement is reconciled per batch of concurrently created rules — at a reduced parallelism each rule waits out a full reconciliation cycle on its own, which can turn an apply that would finish in minutes into one that takes hours.

Note that rate limiting is distinct from the tenant-wide write lock. An `HTTP 409` response reporting `EDIT_LOCK_NOT_AVAILABLE` or `Failed during enter Org barrier` indicates that another session is modifying the configuration, not that a rate limit was exceeded. See [Running Several Configurations Against One Tenant](#running-several-configurations-against-one-tenant).

## Argument Reference - OneAPI

Before starting with this Terraform provider you must create an API Client in the Zscaler Identity Service portal [Zidentity](https://help.zscaler.com/zidentity/what-zidentity) or have create an API key via the legacy method.

* `client_id` - (Required) This is the client ID for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_ID` environment variable.

* `client_secret` - (Optional) This is the client secret for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_SECRET` environment variable. `client_secret` conflicts with `private_key`.

* `private_key` - (Optional) This is the private key for obtaining the API token (can be represented by a filepath, or the key itself). It can also be sourced from the `ZSCALER_PRIVATE_KEY` environment variable. `private_key` conflicts with `client_secret`. The format of the PK is PKCS#1 unencrypted (header starts with `-----BEGIN RSA PRIVATE KEY-----` or PKCS#8 unencrypted (header starts with `-----BEGIN PRIVATE KEY-----`).

* `vanity_domain` - (Optional) This refers to the domain name used by your organization. It can also be sourced from the `ZSCALER_VANITY_DOMAIN`.

* `zscaler_cloud` - (Optional) This refers to Zscaler cloud name where API calls will be directed to i.e `beta`. It can also be sourced from the `ZSCALER_CLOUD`.

**NOTE**: Authentication to the Sandbox service is idependent from authentication to OneAPI or the Legacy API framework and can be set and used in standalone mode.

* `http_proxy` - (Optional) This is a custom URL endpoint that can be used for unit testing or local caching proxies. Can also be sourced from the `ZSCALER_HTTP_PROXY` environment variable.

~> **Rate limiting requires no configuration.** The provider enforces the documented Zscaler API rate limits on the client side and, when a limit is reached, waits out the `Retry-After` interval returned by the API before retrying automatically. `max_retries` and `request_timeout` are the only tuning attributes the provider supports: they bound how long a failing request is retried and how long a single request may take. Neither changes request throughput, so lowering them cannot make `plan` or `apply` faster. Environment variables belonging to the Zscaler Go SDK (for example `ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES` or `ZSCALER_CLIENT_REQUEST_TIMEOUT`) are **not** part of the provider's configuration and are ignored.

* `max_retries` - (Optional) Maximum number of times a rate-limited or transiently failing request is retried before the operation fails. Default: `30`. Maximum: `100`. Each retry waits out the interval reported by the API, so raising it costs nothing when rate limits are not being reached; lowering it only makes runs fail sooner under sustained rate limiting.

* `request_timeout` - (Optional) Timeout in seconds for a single HTTP request to the Zscaler API. Default: `0`, which selects the SDK's built-in 60-second timeout. Accepted values: `0`-`300`. Raise it only if individual requests are timing out.

* `skip_credentials_validation` - (Optional) When set to `true`, the provider skips credential validation and does not initialize the API client. Can also be sourced from the `ZSCALER_SKIP_CREDENTIALS_VALIDATION` environment variable. This is intended for configurations where the ZTC provider is declared but every `ztc_*` resource and data source is conditionally disabled (e.g., `count = 0`) — such as multi-environment deployments where Zscaler is not present in every environment. With this flag enabled, `terraform plan`/`apply` succeeds with a warning even when no credentials are supplied; any resource or data source that does attempt an API call fails with an explanatory error. Default: `false`.

* `username` - (Optional) Administrator account used when authenticating to the legacy Zscaler API framework. Can also be sourced from the `ZTC_USERNAME` environment variable.

* `password` - (Optional) Administrator password used when authenticating to the legacy Zscaler API framework. Can also be sourced from the `ZTC_PASSWORD` environment variable.

* `api_key` - (Optional) API key found in the Zero Trust Cloud  portal `Administration > API Key Mangement > API Key Management`. Can also be sourced from the `ZTC_API_KEY` environment variable. Ensure you have the following SKU enabled `Z_API`

* `ztc_cloud` - (Optional) This refers to the Zscaler cloud name where api calls will be forward to. Can also be sourced from the `ZTC_CLOUD` environment variable.
Currently the following cloud names are supported:
  * `zscaler`
  * `zscalerone`
  * `zscalertwo`
  * `zscalerthree`
  * `zscloud`
  * `zscalerbeta`
  * `zscalergov`
  * `zscalerten`
  * `zspreview`

* `use_legacy_client` - (Optional) This parameter is required when using the legacy API framework. Can also be sourced from the `ZSCALER_USE_LEGACY_CLIENT` environment variable.

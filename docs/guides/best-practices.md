---
page_title: "Best Practices"
description: "Recommended practices for managing Zscaler Zero Trust Cloud (ZTC) with the Zscaler Terraform provider"
---

# Best Practices

This guide collects the recommended practices for managing Zscaler Zero Trust Cloud (ZTC) with Terraform: the general Terraform practices as they apply to this provider, followed by the practices specific to ZTC — activation, OneAPI-only resources, forwarding-rule order and rank, and cloud-orchestrated objects.

-> **Disclaimer:** Please refer to our [General Support Statement](support.md) before proceeding with the use of this provider. For specific errors, see the [Troubleshooting Guide](troubleshooting.md).

## General practices

Terraform practice that applies to the ZTC provider as it does to the other Zscaler providers. ZTC-specific behaviour — activation, forwarding-rule order and rank, and cloud-orchestrated objects — is in [ZTC-specific practices](#ztc-specific-practices).

### Authentication model

- **DO** use OneAPI (Zidentity OAuth2) when the tenant has migrated — `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET` (or `ZSCALER_PRIVATE_KEY`) and `ZSCALER_VANITY_DOMAIN`
- **DO** use the legacy framework (`use_legacy_client = true`) for tenants not yet migrated to Zidentity
- **DO** set `zscaler_cloud` only to target a non-production Zidentity environment (`beta`); omit it for commercial production
- **DON'T** put a legacy cloud name (`zscaler`, `zscalergov` …) in `zscaler_cloud` — those belong in `ztc_cloud` on the legacy path
- **DON'T** leave both credential sets in one job — the path is chosen **only** by `use_legacy_client` / `ZSCALER_USE_LEGACY_CLIENT`, and the other set is silently ignored
- **DON'T** set both `client_secret` and `private_key` — the argument reference documents them as conflicting
- **DON'T** hard-code credentials in the provider block — the docs warn it "risks secret leakage should this file be committed to public version control"

#### Determining the auth path

The provider decides the path from one input: `use_legacy_client`, falling back to `ZSCALER_USE_LEGACY_CLIENT`. When it is false, `ZTC_USERNAME` / `ZTC_PASSWORD` / `ZTC_API_KEY` are never consulted, however correct they are — which is why a job carrying both namespaces fails with a confusing 401 rather than an explicit error.

API access is granted by role. Create an API Role in ZTC (*Administration > Role Management > Add API Role*), then, if the role does not appear in Zidentity, use *Sync Now* under *Integration > API Resources*.

**Do — empty provider blocks; credentials from the environment**

```hcl
provider "ztc" {}   # OneAPI: ZSCALER_CLIENT_ID, ZSCALER_CLIENT_SECRET, ZSCALER_VANITY_DOMAIN
# Legacy: ZTC_USERNAME, ZTC_PASSWORD, ZTC_API_KEY, ZTC_CLOUD, ZSCALER_USE_LEGACY_CLIENT=true
```

**Don't — both namespaces in one job**

```sh
export ZSCALER_CLIENT_ID=...       # OneAPI set
export ZTC_USERNAME=...            # legacy set — silently ignored
export ZSCALER_CLOUD=zscalergov    # legacy cloud name — rejected on OneAPI
```

### Credential handling

- **DO** let the provider read credentials from environment variables, injected from a secret store at job start (Vault, AWS Secrets Manager, Azure Key Vault — reached through your CI's OIDC)
- **DO** keep state in an encrypted, access-controlled remote backend — secrets you set on resources are stored in state in plain text
- **DO** protect saved plan files (`tfplan`) like state — they contain the values of the variables used to produce them
- **DO** prefer `private_key` over `client_secret` where the customer can manage keys
- **DON'T** count on write-only (`*_wo`) arguments or ephemeral resources — Terraform 1.11 supports them, but the provider does not expose any today
- **DON'T** commit `terraform.tfvars` with secrets — commit `terraform.tfvars.example`
- **DON'T** pass secrets with `-var "client_secret=…"` — it lands in shell history and CI logs
- **DON'T** output secrets at all, even marked `sensitive`
- **DON'T** leave `TF_LOG=DEBUG` on in routine CI — debug logs contain request and response bodies

#### Two kinds of secret

**Provider credentials** — client ID, client secret, private key, and the legacy credentials. Provider configuration is never written to Terraform state. When the provider reads these from the environment, as the provider documentation recommends, they never appear in state, in plan output, or in code.

**Secrets that are resource attributes.** Anything a resource accepts as a secret is stored in state in plain text. Marking variables `sensitive = true` only redacts CLI output; it is not encryption. The control is the backend: encryption at rest, tight IAM, and state locking.

**Don't**

```hcl
variable "api_key" {
  default = "9f8a7b6c5d4e"     # in git forever
}

output "api_key" {
  value     = var.api_key
  sensitive = true              # redacted on screen, still readable
}
```

**Do**

```hcl
provider "ztc" {}   # credentials from ZSCALER_* env vars — never in state

variable "psk" {
  type      = string
  sensitive = true   # from TF_VAR_psk; the value still lands in state
}                    # -> encrypt and restrict the backend

# Output only non-sensitive identifiers
output "rule_id" {
  value = ztc_traffic_forwarding_rule.direct_branch_a.id
}
```

-> **Two operational habits** Add `trivy config --severity HIGH,CRITICAL --exit-code 1 .` to CI to catch secrets in tfvars and insecure backend configuration before merge.

Rotate static credentials every 90 days at most.

### Resource addressing: for_each, not count

- Terraform tracks every resource by an **address**. `count` addresses by position, `for_each` addresses by key
- `count` → `rules[0]`, `rules[1]`, `rules[2]` … **the index is the identity**
- Remove the 2nd of 20 items and every later element shifts down one position, so all 18 addresses now point at different objects — Terraform reads that as **18 destroy/recreates**, not one delete
- `for_each` → `rules["block_gambling"]` — the key is the identity, so removing one key touches **only** that object
- **DO** use `for_each` with a map keyed by **purpose**; use `count` only for on/off (`count = var.enabled ? 1 : 0`)

Terraform does not identify a resource by its name or its content — it identifies it by an address in state. With `count`, that address contains the list index, so **the position in the list becomes the object's identity**. Delete the middle entry from the input list and Terraform sees that every later index now holds a different object, so it destroys and recreates them. One deletion becomes N changes. With `for_each` the address carries a stable key, so removing `"streaming"` removes exactly `rules["streaming"]` and nothing else is re-evaluated.

The cost is higher on a ZTC tenant than in a generic Terraform context. Every needless destroy/recreate is an API write, each recreated forwarding rule must be placed back into its order, and in the meantime the rule does not exist. Key by purpose, never by number.

**Don't**

```hcl
resource "ztc_ip_destination_groups" "groups" {
  count = length(var.groups)          # address = groups[0], groups[1] …
  name  = var.groups[count.index].name
  # …
}
```

**Do**

```hcl
resource "ztc_ip_destination_groups" "groups" {
  for_each = var.groups                 # keyed by purpose
  name     = each.key
  # …
}
```

### Drift and ignore_changes

- **Rare, but it happens:** an attribute you set is missing — or returned empty — in the API response
- Terraform stores what the API returns, so the value never lands in state and every plan shows a change
- **This is an API-side issue, not a provider bug** — the provider can only record what it is given
- Common shapes: booleans omitted when false, API-injected defaults, write-only secrets never returned, list ordering, and value normalisation (ZPA lowercases `domain_names`)
- **DO** check the provider CHANGELOG and upgrade first — many drift cases are fixed in a later release
- **DON'T** reach for `ignore_changes` as the *first* move, and never `= all` — it hides genuine drift too

#### Two available responses, in order of preference

1. **Raise a bug with the respective engineering team** so it is fixed at the API backend. This is the real fix.
2. **Temporarily** add `lifecycle { ignore_changes = [<field>] }` scoped to that field until the root cause is resolved.

Both the time limit and the scoping matter. A narrowly scoped suppression with a comment and a ticket number is a documented workaround that someone can find and remove later. `ignore_changes = all` is how a configuration silently diverges from the tenant for months, and it will mask the *next*, unrelated problem too.

Write-only fields are a special case and not a defect: passwords, pre-shared keys and API keys are never returned by design. After an import the secret is not recoverable from the API and must be re-supplied through a variable.

**Don't**

```hcl
# Blanket suppression hides genuine drift as well as the bug
lifecycle {
  ignore_changes = all
}
```

**Do**

```hcl
# Narrowly scoped to the affected field, with a comment and a ticket
# reference. Remove once the API fix ships.
lifecycle {
  ignore_changes = [description]   # SUP-XXXX: absent from GET
}
```

### Parallelism and provider tuning

- **DO** run Terraform with its default parallelism — lowering it slows rule placement and stretches runs
- `-parallelism` is **global to the whole run**. There is no per-resource, per-endpoint or per-provider setting, so lowering it to help one resource type slows down **everything else** in the same configuration
- **Why ZTC is hurt most:** ZTC orders forwarding rules server-side, and the provider reconciles rule placement per **batch** of concurrently created rules. Starve the batch and you pay the reconciliation cost per rule instead of per batch
- The provider **argument** `parallelism` is still honoured in v0.2.0 (an internal request semaphore, default 1) — leave it unset; it is deprecated in the next release
- Seeing lock contention (`EDIT_LOCK_NOT_AVAILABLE`)? The lever is fewer concurrent **runs**, not slower ones — see [Execution and concurrency](#execution-and-concurrency)

#### Why lowering parallelism backfires

Forwarding-rule order is maintained server-side, so after creating rules the provider has to move each rule into its declared position. That reconciliation runs on a cycle shared by all rules created concurrently; at `-parallelism=1` each rule sits through a cycle by itself, and a longer run cannot be offset by raising a session timeout, because ZTC has no such setting.

**Don't**

```sh
terraform apply -parallelism=1   # ~360 rules: 6+ hours

provider "ztc" {
  parallelism = 4                # not a speed-up; deprecated next release
}
```

#### Provider tuning arguments: leave them at the defaults

`max_retries` and `request_timeout` bound how long a failing request is retried and how long one request may take; neither changes throughput. The published v0.2.0 argument reference gives defaults that do not match the code — the table quotes the code.

| Provider | `max_retries` default | `request_timeout` default (range) | Deprecated and ignored |
|---|---|---|---|
| **ZTC** | 30 (maximum 100) | 0 → SDK 60 s (0–300) | None yet — `parallelism` still active in v0.2.0 |

A `request_timeout` of `0` selects the SDK's built-in 60-second timeout; it does **not** disable the timeout.

#### Conditionally disabled providers: skip_credentials_validation

`skip_credentials_validation` — used by the other providers for environments where every resource is disabled — is implemented in the ZTC repository but not in a published release. Until it ships, a declared `ztc` provider needs credentials in every environment.

If a plan is slow on a very large configuration, the cause is the number of resources refreshed at the documented rate limits, not provider settings. Routine changes can plan with `-refresh=false` (in HCP Terraform / TFE, set `TF_CLI_ARGS_plan="-refresh=false"` on the workspace) as long as a scheduled full-refresh plan still detects out-of-band changes; configurations that keep growing should be split. This is the procedure in the ZPA troubleshooting guide (*Slow terraform plan with large configurations*).

### State organisation

- How you lay out state and which backend you use is **your call** — any remote backend works
- **DO** split on whatever boundary matches your teams — by product, by rule type, by environment, by approval group
- **DO** use `terraform_remote_state` (or equivalent) to share IDs across states rather than merging them
- **DON'T** import the same object into two states — every object needs exactly one owner

#### ZTC constraints on layout

| Provider | Constraint on state layout and execution | Section |
|---|---|---|
| **ZTC** | All rules of one type in one configuration · one apply at a time per tenant · one activation per tenant, after the last apply | [Forwarding rules: order, rank and predefined rules](#forwarding-rules-order-rank-and-predefined-rules), [Execution and concurrency](#execution-and-concurrency) |

This section is deliberately light on prescription. State layout and backend choice are architecture decisions for the operator. What *is* ours to say is the constraint above, and it is a consequence of API behaviour rather than an opinion.

None of it argues for one big state: splitting ZIA firewall rules from URL filtering rules is fine, because their orderings are independent. The execution constraints are the ones people get wrong once they have more than one state — splitting state does not split the tenant.

**One workable layout among many**

```text
infrastructure/ztc/
├── prod/forwarding/       # all traffic forwarding rules of each type
└── prod/groups/           # IP destination / source groups, network services
```

### Modules

- **DO** keep `main.tf`, `variables.tf`, `outputs.tf`, `versions.tf`, `README.md` in every module
- **DO** declare `required_providers` in **every** module, and keep the `version` constraint at the environment (root) level only — the provider troubleshooting guides recommend exactly this to avoid *registry does not have a provider named hashicorp/…*
- **DO** configure providers **only in the root module**
- **DO** give each module a single clear **purpose** — purpose is the boundary, not product
- **DO** declare multiple providers in `required_providers` when the *feature* genuinely spans products
- **DON'T** put a `provider` block in a reusable module
- **DON'T** build a `zscaler_everything` god module bundling unrelated products
- **DON'T** hard-code IDs like `"72058304855457833"` — use data sources

#### Purpose as the module boundary

The blunt formulation — "never mix ZIA and ZPA in one module" — is incorrect. The real boundary is **purpose**: "can I describe what this module does in one sentence without using the word *and*?" A module implementing one feature that spans two products — IP Source Anchoring (a ZIA gateway defined by a ZPA server group and segment), or ZTC forwarding rules that take `src_workload_groups` IDs from the ZIA provider — needs both providers. Even then it declares both only in `required_providers` and configures neither; the root module owns authentication, which is what makes aliasing and multi-tenant composition work.

### Version pinning

- **DO** pin to the **minor**: `ztc = "~> 0.2.0"` — the provider is pre-1.0
- **DO** set `required_version` from the features you use: `>= 1.3` minimum (stable `optional()` defaults); `>= 1.7` if you use `removed {}` or mock providers. Every provider CHANGELOG states support for Terraform **v1.x**
- **DO** commit `.terraform.lock.hcl` — the troubleshooting guides show the *Failed to install provider* checksum error you get without it — and upgrade in a dedicated PR tested in non-prod first
- **DO** read the provider CHANGELOG before filing a bug
- **DON'T** omit version constraints
- **DON'T** use `~> 0.1` on a pre-1.0 provider — it accepts 0.2, 0.3 and so on
- **DON'T** run `init -upgrade` inside a feature PR

| Provider | Latest release (verified 15 Oct 2026) | Recommended constraint |
|---|---|---|
| `zscaler/ztc` | v0.2.0 (29 Jun 2026) | `~> 0.2.0` |

The provider is pre-1.0: under semantic versioning a minor bump may change schemas. Terraform reads `~> 0.1` as "any 0.x from 0.1 up"; `~> 0.2.0` accepts only 0.2.x patch releases. Use v0.2.0 (at minimum v0.1.8) if one apply manages more than one rule type ([Forwarding rules: order, rank and predefined rules](#forwarding-rules-order-rank-and-predefined-rules)).

**Do**

```sh
terraform {
  required_version = ">= 1.7.0, < 2.0.0"   # 1.7+: removed {}, mock providers
  required_providers {
    ztc = { source = "zscaler/ztc", version = "~> 0.2.0" }   # pre-1.0: pin the minor
  }
}
```

### Naming, variables and outputs

- All lowercase; underscores in Terraform identifiers, hyphens in names that appear in the Zscaler portal
- Names communicate purpose: `internal_hr_portal`, not `app1` / `a` / `test`
- Carry environment context in a local rather than repeating it: `prefix = "${var.environment}-${var.project}"`
- Separate required from optional variables, and give every variable a `description` and a `type`
- Use `validation` blocks to fail at plan time rather than at the API — particularly for `order`, which the ZTC provider does not check at plan time
- Output identifiers, never secrets ([Credential handling](#credential-handling))

### CI/CD

1. **Validate** — `terraform fmt -check`, `validate`, `tflint`
2. **Scan** — `trivy config`, `checkov`
3. **Plan** — `terraform plan -out=tfplan`, save the artifact
4. **Apply** — `terraform apply tfplan` (the reviewed artifact)
5. **Activate** — `ztcActivator` as its own stage, once, after the last apply

- **DON'T** re-run plan inside the apply job — you'd apply something nobody reviewed
- **DO** keep host-cloud (backend) auth and Zscaler API auth as separate steps
- **DO** serialise applies with a concurrency group keyed on the **tenant** — one tenant-wide write lock

The saved-artifact rule is the one most often bypassed, usually to work around an expired plan. Re-planning in the apply job silently discards human review — it is the CI equivalent of `--auto-approve`.

Authenticating to S3 or Azure Blob for state is a completely different concern from authenticating to the Zscaler API. Keep the two credential paths visibly separate.

With several states against one tenant, activation is one job for the whole pipeline, not one per state. It runs after the last apply, and the concurrency group that serialises the applies is keyed on the tenant rather than the repo.

**With several states against one tenant**

```text
plan (parallel — reads are safe) ──▶ apply (serialised on tenant) ──▶ activate (once)
forwarding · groups                    one at a time                 ztcActivator
```

On HCP Terraform or Terraform Enterprise, routine plans can skip the refresh with the workspace variable `TF_CLI_ARGS_plan="-refresh=false"`. Keep a scheduled full-refresh plan so changes made outside Terraform are still detected.

-> **Two scheduling habits** Run drift detection on a **daily schedule**, not per-commit.

Never run integration tests against production from PR CI.

### Testing

- **DO** run static analysis (`fmt`, `validate`, `tflint`) on every PR
- **DO** run `terraform plan` against a **non-prod tenant** on every PR
- **DO** use `terraform test` (1.6+) for variable validation
- **DO** tag test resources and run a sweeper to clean up
- **DON'T** expect mock providers to catch API behaviour — no rule-order conflicts, no `DUPLICATE_ITEM`, no activation dry-run
- **DON'T** expect plan to catch ordering conflicts — ZTC does not validate `order` at plan time at all
- **DON'T** test in production

Mock providers (Terraform 1.7+) validate *module wiring*, not Zscaler semantics. They are valuable and cheap, but a green mock suite says nothing about whether the API accepts the configuration. Pair them with a sandbox tenant.

~> **Activation in test teardown** Test teardown **must include activation** — `ztcActivator`, run after `terraform destroy` as the activator guide shows. Otherwise the deletes stay as drafts.

Pass credentials to the activator through environment variables, not CLI arguments — the guide calls CLI arguments insecure.

### Import and brownfield adoption

- **DO** use **zscaler-terraformer** for brownfield adoption — `import` writes the HCL *and* imports into state; `generate` writes HCL only
- **DO** adopt **one resource type per PR**, and confirm an empty `terraform plan` before the next
- **DO** import individual objects by ID (many resources also accept the name — check the resource page)
- **DO** keep the default parallelism, even for a large initial import
- **DO** re-supply write-only secrets after import — the API won't return them
- **DO** use `removed {}` (Terraform 1.7+) to stop managing something without destroying it
- **DON'T** hand-write HCL for a tenant that already exists — that is what the tool is for
- **DON'T** use `terraform state rm` on a custom object — it orphans the API object and the next apply hits `DUPLICATE_ITEM`

zscaler-terraformer is a supported Zscaler CLI that authenticates with the same credentials as the providers, reads a live tenant over the API, and emits the Terraform configuration and the import into state, so the first plan after an import should come back clean. Scope it with `--resources`: the product name for a full sweep, or a single resource type when onboarding incrementally — the pattern to recommend. `--prefix` namespaces the generated names, `--validate` runs `terraform init` and `terraform validate` over the output, and `--collect-logs` captures SDK debug logs for a support case.

Generated HCL is correct but not idiomatic: names come from the API, resources are flat rather than grouped with `for_each`, and there are no variables. Refactor it before it becomes the long-term source of truth.

**Do — brownfield: don't hand-write it**

```sh
zscaler-terraformer --supported-resources="ztc"                 # list what is covered
zscaler-terraformer import --resources "ztc"
zscaler-terraformer --prefix "prod" --validate import --resources "ztc"
```

-> **Coverage and caveats** The tool's README introduction still mentions only ZPA and ZIA, but it ships a ZTC provider package and a *ZTC Supported Resources* table. Confirm that the resource types you need are in that table; cloud-orchestrated objects remain data sources either way ([Scope and resource model](#scope-and-resource-model)).

It cannot recover write-only values — passwords, pre-shared keys, API keys — because the API never returns them. Re-supply those from your secret store after the import.

Supply credentials as environment variables rather than CLI flags so they do not land in shell history.

`terraform state rm` is the most frequently misused command in this context. It appears to offer a clean exit from a bad state, but it leaves the object live in the tenant and untracked, and the next apply tries to create it again. `removed {}` expresses the same intent safely and is visible in code review.

**Don't**

```sh
terraform state rm ztc_ip_destination_groups.branch
# Object still exists in the tenant. Next apply -> DUPLICATE_ITEM or a duplicate.
```

**Do**

```hcl
removed {
  from = ztc_ip_destination_groups.branch
  lifecycle { destroy = false }
}
```

## ZTC-specific practices

ZTC has two constraints in common with ZIA — changes are drafts until activated, and writes serialise behind a tenant-wide lock — but it is a separate provider with its own documentation, its own activator and its own differences: no adjustable session timeout, resources that work only over OneAPI, forwarding rules ordered by both order and admin rank, objects created by cloud orchestration, and a pre-1.0 schema. Everything here is built from the ZTC provider documentation (docs/index.md, the ztc-activator-overview and troubleshooting guides, the resource pages), the CHANGELOG and the source at v0.2.0; where the published docs and the code disagree, it quotes the code.

-> **Naming** ZTC was formerly Cloud Branch Connector. The rebrand affects branding only: resource names keep the `ztc_` prefix and the API path remains `/cloud-branch-connector`.

### Scope and resource model

The zscaler/ztc provider registers **16 resources and 20 data sources** (ztc/provider.go at v0.2.0). Four objects are exposed as data sources with no resource counterpart, because another system owns them.

| Data source | Why it is read-only |
|---|---|
| `ztc_edge_connector_group` | Created by AWS, Azure or GCP cloud connector orchestration when a connector instance starts |
| `ztc_location_management` | Auto-provisioned per VPC or VNet by the cloud integration |
| `ztc_supported_regions` | Cloud-provider region catalog |
| `ztc_workload_groups` | Defined in ZIA and referenced from ZTC |

#### Cloud-orchestrated objects

This is the most common ZTC-specific error. Edge connector groups and locations are created automatically when a cloud connector instance starts. A resource block for one either fails because no such resource exists or collides with the orchestrated object and returns `DUPLICATE_ITEM`. Read them with data sources.

**Don't**

```hcl
resource "ztc_location_management" "aws_vpc" {             # collides with the
  name = "AWS-CAN-ca-central-1-vpc-05c7f364cf47c2b93"   # orchestrated object
}
```

**Do**

```hcl
data "ztc_location_management" "aws_vpc" {
  name = "AWS-CAN-ca-central-1-vpc-05c7f364cf47c2b93"
}
resource "ztc_traffic_forwarding_rule" "this" {
  locations { id = [data.ztc_location_management.aws_vpc.id] }
  # …
}
```

### Authentication and OneAPI-only resources

- **DO** confirm the auth path before writing the provider block
- **DO** use legacy (`use_legacy_client = true`, `ztc_cloud = "zscalergov"` or `"zscalerten"`) for government ZTC tenants
- **DON'T** assume FedRAMP over OneAPI works because it does on ZIA and ZPA — it is not in a published ZTC release
- **DON'T** set `zscaler_cloud = "zscaler"` on OneAPI — `zscaler` is a legacy cloud name
- **DON'T** mix `ZSCALER_*` and `ZTC_*` variables in one job
- **DON'T** manage `ztc_traffic_forwarding_dns_rule` or `ztc_traffic_forwarding_log_rule` over legacy — both pages say "Resource available only via OneAPI"
- **DO** read *Unauthorized access to private resource* as "this resource needs OneAPI" — the troubleshooting guide: "the resource is NOT accessible via the Legacy API and only via OneAPI"
- **DO** get `src_workload_groups` IDs from the **ZIA** provider — the traffic forwarding rule examples: "you must leverage the ZIA Terraform Provider at the moment"
- **DO** use ISO-3166 Alpha-2 for `countries` on `ztc_ip_destination_groups`

**OneAPI — Zidentity tenants**

```sh
terraform {
  required_providers {
    ztc = { source = "zscaler/ztc", version = "~> 0.2.0" }   # pre-1.0: pin the minor
  }
}
provider "ztc" {}
# ZSCALER_CLIENT_ID / ZSCALER_CLIENT_SECRET (or ZSCALER_PRIVATE_KEY) / ZSCALER_VANITY_DOMAIN
# ZSCALER_CLOUD   (optional — non-production only, e.g. "beta")
```

**Legacy — pre-Zidentity tenants, and all government tenants**

```hcl
provider "ztc" {
  use_legacy_client = true
  ztc_cloud         = "zscalergov"
}
# ZTC_USERNAME, ZTC_PASSWORD, ZTC_API_KEY, ZSCALER_USE_LEGACY_CLIENT=true

# Workload groups come from ZIA, referenced in ZTC
data "zia_workload_groups" "apps" { name = "Apps" }
```

- **DO** inject credentials as environment variables from a secret store
- **DON'T** commit `client_secret`, `private_key`, `password` or `api_key` in `.tfvars`

[Credential handling](#credential-handling) applies without modification. The legacy ZTC path uses a username, password and API key rather than a client ID and secret, so three values require protection rather than two.

### Activation

**The docs: "The ZTC platform requires every configuration to be activated."**

- **DON'T** assume `terraform apply` = enforced — every create, update and delete is a draft until activated
- **DO** run `ztcActivator` once, out of band, after the last `apply` **or** `destroy` — built with `make ztcActivator` from the provider source (you may need `sudo` to install it into `/usr/local/bin`)
- **DON'T** pass credentials to the activator on the command line — the guide: "Providing authentication credentials via CLI argument is insecure and is not recommended"
- **DO** use `ztc_activation_status` in HCL only when a single flat state owns all its resources, with an exhaustive `depends_on` — the resource is registered in v0.2.0 but has no published documentation page yet
- **DON'T** activate once per state when several states target one tenant — activation is tenant-wide
- Pure data-source configurations — no `resource "ztc_…"` blocks — need no activation

Activation is tenant-wide: one call publishes every pending change, so no configuration needs more than one activation per run. Pending changes are also activated when an API session ends, so a long apply can publish a half-written configuration part-way through — and unlike ZIA, ZTC exposes no session-timeout setting to lengthen that window. Keeping runs short, by splitting state rather than lowering parallelism, is the only control.

**Do**

```sh
terraform init && terraform apply && ztcActivator
terraform destroy && ztcActivator    # deletes are drafts too
```

**Also fine — flat state only**

```hcl
resource "ztc_activation_status" "this" {
  status     = "ACTIVE"
  depends_on = [ztc_traffic_forwarding_rule.direct_branch_a]
}
```

### Forwarding rules: order, rank and predefined rules

Applies to `ztc_traffic_forwarding_rule`, `ztc_traffic_forwarding_dns_rule` and `ztc_traffic_forwarding_log_rule`.

- **DO** keep orders contiguous from 1, with all rules of one type in **one** configuration
- **DON'T** rely on plan to catch a bad `order` — ZTC does **not** validate it at plan time; the API rejects it at apply
- **DO** respect **admin rank** — the API refuses `Rule with rank N is not allowed at order M`; the documented examples use `rank = 7`
- **DON'T** destroy the predefined forwarding rules — the provider refuses with `predefined rule '…' cannot be deleted` (*ZPA Forwarding Rule*, *Direct rule for Zscaler Cloud Endpoints*, *Direct rule for WAN / LAN Destinations Group*, *Client Connector to ZPA*, *ZPA Pool For Stray Traffic*). Release them with `removed {}`
- **DO** renumber every rule below a deleted one, and don't leave a gap
- **DO** run v0.1.7 or later when one apply manages several rule types (reorder race fix, PR #24)

Forwarding-rule order is a contiguous server-side sequence: deleting any rule other than the last shifts every rule below it up one. All rules of one type belong in one configuration, because two configurations computing positions for the same list will fight. When the order you ask for conflicts with a rank, the provider surfaces the API error together with the current order-versus-rank layout. If you build orders from input, add a variable `validation` block — the provider will not catch a bad value before apply.

**Don't**

```hcl
resource "ztc_traffic_forwarding_rule" "a" {
  order = 0      # not caught at plan — fails at apply
}
```

**Do**

```hcl
resource "ztc_traffic_forwarding_rule" "direct_branch_a" {
  name           = "DIRECT_branch_a"
  order          = 1          # contiguous from 1
  rank           = 7          # default admin rank
  state          = "ENABLED"
  type           = "EC_RDR"
  forward_method = "DIRECT"
}
```

### Execution and concurrency

- **DO** serialise apply with a concurrency group keyed on the tenant; keep plan parallel — reads do not take the write lock
- **DO** run with Terraform's default parallelism
- **DON'T** lower `-parallelism` — rule placement is reconciled per batch, so the cost per rule rises ([Parallelism and provider tuning](#parallelism-and-provider-tuning))
- **DON'T** raise the provider argument `parallelism` as a speed-up — in v0.2.0 it drives an internal request semaphore (default 1), and it is deprecated in the next release

ZTC serialises writes behind a single tenant-wide lock. Several states applying concurrently contend for it and produce `EDIT_LOCK_NOT_AVAILABLE`; splitting state does not split the tenant. Provider defaults in v0.2.0 are `max_retries` 30 and `request_timeout` 0 (the SDK's 60 seconds) — the published argument reference still says 5 and "no limit".

**Do — serialise on the tenant**

```yaml
# GitHub Actions
concurrency:
  group: ztc-apply-${{ inputs.tenant }}
  cancel-in-progress: false
```

### State, versions and import

- **DO** keep one Terraform configuration per ZTC tenant where practical — activation is tenant-wide
- **DO** pin to the **minor**: `version = "~> 0.2.0"` (latest tagged release v0.2.0, 29 June 2026)
- **DON'T** use `~> 0.1` or `>= 0.1` on a pre-1.0 provider — both accept the next minor
- **DO** use the five-stage pipeline ([CI/CD](#cicd)) with `ztcActivator` as the activation stage, and include activation in test teardown
- **DO** use zscaler-terraformer for brownfield ZTC tenants — list coverage with `--supported-resources="ztc"`

### Troubleshooting reference

| Error | Likely cause | Fix |
|---|---|---|
| Apply succeeds, configuration not in effect | No activation step | Run `ztcActivator` after apply |
| `Unauthorized access to private resource` | OneAPI-only resource (DNS or log forwarding rule) on legacy | Use OneAPI |
| `DUPLICATE_ITEM` on a location | A resource block for a cloud-orchestrated object | Use `data "ztc_location_management"` |
| `EDIT_LOCK_NOT_AVAILABLE` | Concurrent applies against one tenant | Serialise applies on the tenant |
| `order = 0` rejected at apply | Zero or negative order value | Set `order >= 1` |
| `Rule with rank N is not allowed at order M` | Requested order conflicts with an admin rank | Reorder respecting rank |
| `predefined rule '…' cannot be deleted` | Destroy on a predefined forwarding rule | `removed {}` with `destroy = false` |
| Rules drift after a delete | Survivors not renumbered in HCL | Decrement every order below the deleted rule |
| Rule reorder race across rule types | Provider below v0.1.7 | Upgrade to v0.2.0 |
| Slow applies with no errors | `parallelism` provider argument (default 1), or lowered `-parallelism` | Leave both at the defaults |
| GOV cloud rejected on OneAPI | FedRAMP over OneAPI not in a published ZTC release | `use_legacy_client = true` with `ZTC_CLOUD` |
| Schema changed after a minor upgrade | Pre-1.0 provider with a loose constraint | Pin `~> 0.2.0` (or exact) |

# Terraform Provider ZTC — Claude Code Guidelines

This file provides project-specific guidance for the `terraform-provider-ztc` Terraform provider. Follow these conventions when creating, modifying, or reviewing code.

## Project Overview

This is a Terraform provider for **Zscaler Zero Trust Cloud (ZTC)**—the Cloud & Branch Connector control plane historically exposed under ZTW APIs. The provider wraps the **ZTW REST APIs** via the shared Go SDK (`github.com/zscaler/zscaler-sdk-go/v3`). Resources and data sources use **`terraform-plugin-sdk/v2`** and live under the `ztc/` package namespace (provider type **`ztc`** on the Registry).

Important distinction from **ZIA**: configuration **activation** is not triggered automatically after every resource mutation. Operators use the dedicated **`ztc_activation_status`** resource and/or the separate **`ztcActivator`** CLI (see below).

---

## Architecture

### API surface (`ztw` in the SDK)

- SDK imports use paths under **`zscaler/ztw/services/…`** (e.g. policy resources, forwarding rules, forwarding gateways, activation, dns gateway—not `zscaler/zia/services/…`).
- The configured `*zscaler.Service` is what every resource passes into SDK calls.

### Client wrapping

```go
zClient := meta.(*Client)
service := zClient.Service
```

For **Create** contexts, prefer the type assertion with a clear diagnostic (pattern used across resources):

```go
zClient, ok := meta.(*Client)
if !ok {
    return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
}
service := zClient.Service
```

### Dual authentication: OneAPI vs legacy ZTW

`ztc/config.go` + `provider.go` support:

1. **OneAPI** (default): OAuth-style credentials — `client_id` / `client_secret` (+ optional `vanity_domain`, **`zscaler_cloud`** optional; env fallbacks **`ZSCALER_*`**).
2. **Legacy ZTW**: session-based — `username` / `password` / `api_key` mapped to **`ZTC_USERNAME`**, **`ZTC_PASSWORD`**, **`ZTC_API_KEY`**, and cloud/base URL via **`ztc_cloud`** / **`ZTC_CLOUD`**. Enable with provider argument **`use_legacy_client`** or env **`ZSCALER_USE_LEGACY_CLIENT=true`**.

`private_key` is an alternate OneAPI credential path (conflicts with `client_secret` in schema).

### Global request semaphore

`provider.go` allocates **`apiSemaphore`** from configured **`parallelism`** (default channel depth 1) to throttle concurrent API usage where bulk operations are not available. Respect existing patterns when adding resources that coordinate many calls.

---

## Project Structure

```
ztc/
  provider.go                    # Provider registration, schema, semaphore init
  config.go                     # Credential loading, OneAPI + legacy client wiring
  client.go                     # Client type(s) surfaced to resources
  common.go                     # Shared schema helpers, reorder engine, flatten/expand utilities
  utils.go                      # General helpers, DetachRuleIDNameExtensions, etc.
  validator.go                  # Validation helpers
  provider_test.go             # Acceptance test globals, sweep registration, PreCheck
  provider_sweeper_test.go      # Sweeper implementations
  resource_ztc_<name>.go       # Resources
  data_source_ztc_<name>.go     # Data sources
  *_test.go                     # Acceptance tests (TF_ACC=1)
  common/
    version.go                   # Provider version string (releases)
    resourcetype/resource_type.go  # Test + sweeper constants for Terraform type names
    testing/method/, testing/variable/  # Test name generators / shared literals
cli/
  ztcActivator.go               # Standalone binary: OneAPI activation (make ztcActivator)
docs/resources/, docs/data-sources/ # Registry-facing docs (terraform-plugin-docs)
docs/guides/                  # Operational guides (release-notes, activator, etc.)
examples/<resource>/          # Referenced snippets
.claude/skills/               # plan-tf-resource, troubleshoot-resource workflows
.cursor/rules/                # Maintainer rules (provider + troubleshooting)
```

---

## Activation Model (critical)

### In Terraform: `ztc_activation_status`

- **Do not** add per-resource “activate after Create/Update/Delete” hooks in ordinary resources. Policy and routing changes typically require an explicit activation step separately from this provider’s resource CRUD.
- The **`ztc_activation_status`** resource calls `activation.UpdateActivationStatus` via the SDK. **Delete is a no-op** (`resourceFuncNoOp`); design matches “fire activation” rather than owning long-lived config objects.

### Out-of-band: `ztcActivator` CLI

- Built with **`make ztcActivator`** (see `GNUmakefile`). Same credential story as the provider: OneAPI env vars (with **`ZSCALER_CLOUD` optional**, matching `cli/ziaActivator.go` / `ztc/config.go`) or legacy ZTW env vars when `ZSCALER_USE_LEGACY_CLIENT=true`.
- Used when teams want **`terraform apply && ztcActivator`** without managing `ztc_activation_status` in HCL. Documented in `docs/guides/ztc-activator-overview.md` and `docs/guides/ztc_activator.md`.

---

## Ordered rule resources (traffic forwarding)

These resources participate in a **shared reorder coordinator** in `ztc/common.go`:

- `ztc_traffic_forwarding_rule`
- `ztc_traffic_forwarding_dns_rule`
- `ztc_traffic_forwarding_log_rule`

### How reordering works (high level)

1. Each rule registers its desired **`OrderRule` (`Order`, `Rank`)** and marks progress with **`markOrderRuleAsDone`** / **`reorderWithBeforeReorder`**.
2. **`reorderAll`** runs on a **ticker** (`reorderTickInterval`, default 30s; tests may shorten it). It sorts registered rules, optionally runs **`beforeReorder`**, then issues **`updateOrder`** per rule when the API-reported count allows.
3. **`waitForReorder`** must be used after rules are marked done and **before** final **Read** so Terraform state reflects post-reorder API order.
4. A per–resource-type **`reorderDone`** channel allows **finishing one reorder goroutine and starting another** when new rules arrive late (avoids races when multiple rule types or batches apply together).

### When changing reorder logic

- Prefer **extending** `reorderAll` / `reorderWithBeforeReorder` / `waitForReorder` rather than per-resource ad-hoc loops.
- Regressions to watch for: applies that **hang** waiting for reorder, **wrong final order** when many rules apply in parallel, or **duplicate PUTs**. Compare behavior with recent **CHANGELOG** / PR notes when debugging.

Rule **`order`** is **meaningful to the API**; keep schemas and validations aligned with how forwarding rules are documented (`order`/`rank`). When in doubt, mirror existing forwarding rule resources and run acceptance tests.

---

## References and deletes: `DetachRuleIDNameExtensions`

Several policy objects (**IP destination groups, IP source groups, network services, etc.**) appear on **traffic forwarding rules**. Before **Delete**, some resources must **detach** references so the API allows removal. Use **`DetachRuleIDNameExtensions`** (`ztc/utils.go`) with the appropriate getter/setter on `forwarding_rules.ForwardingRules`—follow an existing sibling resource (e.g. IP destination groups) when adding similar objects.

---

## New resources and data sources

Use **`.claude/skills/plan-tf-resource/SKILL.md`** for the full checklist. In summary:

| Step | Action |
|------|--------|
| Registration | Add to `ResourcesMap` / `DataSourcesMap` in `provider.go` |
| Constants | Add to `ztc/common/resourcetype/resource_type.go` for tests sweepers |
| Tests | Resource test with create/update/import + shared config including **paired `data`** block; separate data source test with **`TestCheckResourceAttrPair`** |
| Sweeper | `setupSweeper` in `provider_test.go` + implementation in `provider_sweeper_test.go` when the resource deletes remote objects |
| Variables | Shared strings in `ztc/common/testing/variable/variable.go` when reused |
| Docs | `docs/resources/ztc_<name>.md`, `docs/data-sources/ztc_<name>.md`, `examples/ztc_<name>/` |
| Release | Bump `ztc/common/version.go`, **`GNUmakefile`** `build13` lines (all three), **`CHANGELOG.md`**, **`docs/guides/release-notes.md`** |

**Import:** Support import by **numeric ID** and **name** where the SDK exposes `Get` + `GetByName` (standard `ResourceImporter` pattern in existing resources).

**IDs in state:** Terraform **`id`** is usually a string; many resources also expose **`<thing>_id`** or **`gateway_id`** as the integer API id—keep consistency with neighboring resources.

---

## Schema conventions

- SDK fields with **`omitempty`** on **booleans**: prefer **`Optional: true, Computed: true`** to avoid perpetual drift when the API omits `false`.
- **API-defaulted** optional fields: often **`Optional: true, Computed: true`**.
- **`description`**: follow existing **`StateFunc` / `DiffSuppressFunc`** patterns for multiline text where used.
- **Nested blocks** with API-assigned IDs: nested **`id`** **`Computed: true`**; **expand** must preserve IDs on update (`getIntFromNested`).
- Reuse helpers in **`common.go`** / **`utils.go`** / **`validator.go`** instead of duplicating.

---

## Data sources

- Lookup pattern: optional **`id`** (int) and/or **`name`** (string), both often **Computed + Optional**; implement **Get** vs **GetByName** fallbacks like existing data sources.
- Many read-only objects are **cloud-managed**; prefer **data sources** for those per product expectations (see `.cursor/rules` and **ztc-skill** for “data-source-first” guidance).

---

## Troubleshooting workflow

Prefer **`.cursor/rules/troubleshoot-ztc-provider.md`** and **`.claude/skills/troubleshoot-resource/SKILL.md`** for structured diagnosis (drift, 400s, ordering, stale test data).

**Debug logging:**

```bash
TF_LOG=DEBUG ZSCALER_SDK_VERBOSE=true ZSCALER_SDK_LOG=true terraform apply -no-color 2>&1 | tee /tmp/tf-debug.log
```

---

## JMESPath / client-side `search`

The ZTC provider **does not** currently wire JMESPath **`search`** attributes on data sources (unlike the ZIA provider). If that is added later, mirror the ZIA pattern: `zscaler.ContextWithJMESPath`, optional schema attribute, and doc the filterable JSON field names (**camelCase**).

---

## User-facing writing conventions

Applies to **`docs/**/*.md`**, **`CHANGELOG.md`**, schema **`Description`** strings, and Registry text:

1. **Do not** name SDK packages, Go types, or internal helpers in user-facing prose (use neutral product language: “the resource returns…”, “after activation…”).
2. Reserve **“breaking change”** for real compatibility breaks—see **ZIA `CLAUDE.md`** for the same bar; argument renames with preserved capability are **inline notes**, not banners.
3. **Changelog / release-notes** entries lead with **user-visible outcomes**, name affected **`ztc_*`** resources/data sources, and link PRs as **`[PR #NNN](https://github.com/zscaler/terraform-provider-ztc/pull/NNN)`**.

---

## Build, test, and tooling

```bash
# Build (if vendor is stale, use -mod=mod)
go build ./...
go build -mod=mod ./...

# Unit tests (non-acceptance)
make test-unit

# Single acceptance test (requires ZSCALER_CLIENT_ID / ZSCALER_CLIENT_SECRET / ZSCALER_VANITY_DOMAIN, etc.)
TF_ACC=1 go test ./ztc/ -v -run TestAcc<ResourceName>_Basic -timeout 120m

# Sweepers (dangerous — dev accounts only)
go test ./ztc/ -v -sweep=global -timeout 30m
```

---

## Release versioning

Every release **must** update:

1. `ztc/common/version.go`
2. `GNUmakefile` — **all three** version strings under **`build13`**
3. `CHANGELOG.md` — new section at top
4. `docs/guides/release-notes.md` — mirror entry + **`Last updated: v…`** line

---

## Critical rules (ZTC-specific)

1. **No inline activation** in ordinary resource CRUD—use **`ztc_activation_status`** and/or **`ztcActivator`**; keep that separation when adding resources.
2. **Never** ship a new deletable resource without **acceptance tests**, **docs**, **`examples/`**, and normally a **sweeper**.
3. **Register** everything in **`provider.go`** and test type constants in **`resourcetype/resource_type.go`** when tests/sweepers need it.
4. **Import** by **ID** and **name** where the SDK supports both.
5. **Forwarding-rule-attached objects:** use **`DetachRuleIDNameExtensions`** (or equivalent) before delete when the API requires it—copy an existing resource’s pattern.
6. **Forwarding rules:** after changes, respect **`markOrderRuleAsDone`**, **`reorderWithBeforeReorder`**, and **`waitForReorder`** so ordering stays correct under parallel applies.
7. **SDK path** is **`ztw`**, not **`zia`**.
8. Reuse **`common.go` / `utils.go`** helpers—avoid parallel implementations.
9. **User-facing** copy stays free of SDK/internal jargon; keep deep mechanics here and in `.cursor/rules` / `.claude/skills/`.

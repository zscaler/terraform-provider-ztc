---
layout: "zscaler"
page_title: "Release Notes"
description: |-
  The Zscaler Zero Trust Cloud (ZTC) provider Release Notes
---

# ZTC Provider: Release Notes

## USAGE

Track all ZTC Terraform provider's releases. New resources, features, and bug fixes will be tracked here.

---
``Last updated: v0.2.1``

---

## 0.2.1 (November 20, 2026)

### Notes

- Supported Terraform version: **v1.x**

### Features

- [PR #58](https://github.com/zscaler/terraform-provider-ztc/pull/58) - Added the provider attribute `skip_credentials_validation` (env var `ZSCALER_SKIP_CREDENTIALS_VALIDATION`). When enabled, the provider skips credential validation and API client initialization so that configurations where every `ztc_*` resource and data source is conditionally disabled (e.g., `count = 0`) can plan and apply without credentials — e.g., multi-environment deployments where Zscaler is not present in every environment. A warning is emitted at configure time, and any resource or data source that does attempt an API call fails with an explanatory error instead of a panic. See ([ZPA Issue #684](https://github.com/zscaler/terraform-provider-zpa/issues/684)) for the motivating use case.

### Bug Fixes

- [PR #58](https://github.com/zscaler/terraform-provider-ztc/pull/58) - The `ztc_network_services` and `ztc_network_service_groups` data sources, and import of the `ztc_network_services` and `ztc_network_service_groups` resources by name, no longer fail with `INVALID_INPUT_ARGUMENT` `Invalid Page Size 1,000`. The API now rejects a page size of 1,000, which the lookups requested. Both lookups now read the whole list, 127 items per page, so they also find services and groups beyond the first page.
- [PR #58](https://github.com/zscaler/terraform-provider-ztc/pull/58) - Creating, updating or deleting many objects of the same type no longer triggers waves of rejected requests or long waits. Requests now follow the rate limits documented for each method and resource (for most resources, 2 reads and 1 create, update or delete per second, and 1,000 reads and 400 writes per hour). Before, several requests to the same resource were sent at once, the API rejected most of them, and the rejected requests counted toward the hourly limit. Requests that the API rejects no longer count toward that limit.
- [PR #58](https://github.com/zscaler/terraform-provider-ztc/pull/58) - Upgraded to Zscaler-SDK-GO v3.8.62.

### Documentation

- [PR #58](https://github.com/zscaler/terraform-provider-ztc/pull/58) - Corrected the provider argument reference for `max_retries` (the default is `30`, not `5`) and `request_timeout` (the default is `0`, which selects the SDK's built-in 60-second timeout rather than disabling the timeout), and added a note explaining that rate limiting requires no configuration, that these are the only supported tuning attributes, that neither affects request throughput, and that Zscaler Go SDK environment variables such as `ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES` are not part of the provider's configuration and are ignored.

## 0.2.0 (June 29, 2026)

### Notes

- Release date: **(June 29, 2026)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #45](https://github.com/zscaler/terraform-provider-ztc/pull/45) - Fixed `ztc_edge_connector_group` panic with attribute type `ec_instances.0.dns_ip`

## 0.1.9 (May 13, 2026)

### Notes

- Release date: **(May 13, 2026)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #38](https://github.com/zscaler/terraform-provider-ztc/pull/38) - Updated `ztcActivator` OneAPI startup so `ZSCALER_CLOUD` is optional (unset uses the default production cloud).

- [PR #38](https://github.com/zscaler/terraform-provider-ztc/pull/38) - Fixed resource and datasource `ztc_account_groups` to allow for proper simplified configuration.

## 0.1.8 (February 27, 2026)

### Notes

- Release date: **(February 27, 2026)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #27](https://github.com/zscaler/terraform-provider-ztc/pull/27) - Added new resource and data source `ztc_dns_gateway` for managing DNS Gateway configurations.

### Bug Fixes

- [PR #27](https://github.com/zscaler/terraform-provider-ztc/pull/27)  - Fixed `UpdateContext` in all ordered rule resources to defer `ReadContext` until after the reorder mechanism completes, preventing stale order values from being stored in state.


## 0.1.7 (February 11, 2026)

### Notes

- Release date: **(February 11, 2026)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #24](https://github.com/zscaler/terraform-provider-ztc/pull/24) - Fixed rule reorder race condition where the reorder timer could fire before all rules of a type were registered, causing incorrect ordering when multiple rule types are deployed together. The reorder goroutine now runs asynchronously with automatic re-run for late-arriving rules, and all rules wait for the reorder to complete before reading final state.

## 0.1.6 (February 9, 2026)

### Notes

- Release date: **(February 9, 2026)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #23](https://github.com/zscaler/terraform-provider-ztc/pull/23) - Fixed `ztc_traffic_forwarding_rule` READ function to accomodate OneAPI vs Legacy API availability.

## 0.1.5 (February 5, 2026)

### Notes

- Release date: **(February 3, 2026)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #22](https://github.com/zscaler/terraform-provider-ztc/pull/22) - Fixed `ztc_traffic_forwarding_rule` READ function to accomodate OneAPI vs Legacy API availability.


## 0.1.4 (February 3, 2026)

### Notes

- Release date: **(February 3, 2026)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #21](https://github.com/zscaler/terraform-provider-ztc/pull/18) - Upgraded to Zscaler-SDK-GO v3.8.15
- [PR #21](https://github.com/zscaler/terraform-provider-ztc/pull/21) - Fixed `ztc_traffic_forwarding_dns_rule`,  `ztc_traffic_forwarding_rule` and `ztc_traffic_log_forwarding_rule` resource reorder logic due to recent API enforcement changes. Included safeguard to prevent unnecessary reordering when the order is already correct. 

## 0.1.3 (January 19, 2026)

### Notes

- Release date: **(January 19, 2026)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #18](https://github.com/zscaler/terraform-provider-ztc/pull/18) - Upgraded to Zscaler-SDK-GO v3.8.13
[PR #18](https://github.com/zscaler/terraform-provider-ztc/pull/18) - Fixed resource and datasource `ztc_traffic_forwarding_rule` issue to retrieve by ID.

## 0.1.2 (January 13, 2026)

### Notes

- Release date: **(January 13, 2026)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #16](https://github.com/zscaler/terraform-provider-ztc/pull/16) - Fixed Legacy Client instantiation.


## 0.1.1 (December 1, 2025)

### Notes

- Release date: **(December 1, 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #12](https://github.com/zscaler/terraform-provider-ztc/pull/12) - Fixed resources and data sources `ztc_provisioning_url`, `ztc_location_template`.


## 0.1.0 (December 1, 2025) - 🎉Initial Release🎉

### Notes

- Release date: **(December 1, 2025)**
- Supported Terraform version: **v1.x**

This new Provider (ZTC Terraform Provider) enables fully automated provisioning and configuration across several key components in the Cloud & Branch Connector portal.

### NEW - DATA SOURCES

[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_activation_status`` - Activate ZTC Configuration information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_location_template`` - Retrieves location template information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_provisioning_url`` - Retrieves provisioning URL information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_location_management`` - Retrieves location management information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_edge_connector_group`` - Retrieves Edge Connector Group information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_traffic_forwarding_rule`` - Retrieves Traffic Forwarding rule
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_traffic_forwarding_dns_rule`` - Retrieves Traffic Forwarding DNS rule
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_traffic_forwarding_log_rule`` - Retrieves Traffic Forwarding LOG rule
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_forwarding_gateway`` - Retrieves ZIA and LOG Forwarding Gateway
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_dns_forwarding_gateway`` - Retrieves DNS Forwarding Gateway
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_ip_destination_groups`` - Retrieves IP Destination Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_ip_source_groups`` - Retrieves IP Source Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_ip_pool_groups`` - Retrieves IP Pool Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_network_services`` - Retrieves Network Services
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_network_service_groups`` - Retrieves Network Services Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_account_groups`` - Retrieves Account Groups
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_public_cloud_info`` - Retrieves Public Cloud Info
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_supported_regions`` - Retrieves Supported Regions
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Data Source ``ztc_workload_groups`` - Retrieves Workload Groups

### NEW - RESOURCE

[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_activation_status`` - Activate ZTC Configuration information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_location_template`` - Retrieves location template information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_provisioning_url`` - Retrieves provisioning URL information
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_traffic_forwarding_rule`` - Retrieves Traffic Forwarding rule
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_traffic_forwarding_dns_rule`` - Retrieves Traffic Forwarding DNS rule
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_traffic_forwarding_log_rule`` - Retrieves Traffic Forwarding LOG rule
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_forwarding_gateway`` - Retrieves ZIA and LOG Forwarding Gateway
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_dns_forwarding_gateway`` - Retrieves DNS Forwarding Gateway
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_ip_destination_groups`` - Retrieves IP Destination Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_ip_source_groups`` - Retrieves IP Source Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_ip_pool_groups`` - Retrieves IP Pool Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_network_services`` - Retrieves Network Services
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_network_service_groups`` - Retrieves Network Services Group
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_account_groups`` - Retrieves Account Groups
[PR #11](https://github.com/zscaler/terraform-provider-ztc/pull/11) - Resource ``ztc_public_cloud_info`` - Retrieves Public Cloud Info


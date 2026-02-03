# Changelog

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
- [PR #18](https://github.com/zscaler/terraform-provider-ztc/pull/18) - Fixed resource and datasource `ztc_traffic_forwarding_rule` issue to retrieve by ID.

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

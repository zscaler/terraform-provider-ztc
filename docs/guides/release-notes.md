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
``Last updated: v0.1.0``

---

# Changelog

## 0.1.0 (December 1, 2025) - 🎉Initial Release🎉

### Notes

- Release date: **(December 1, 2025)**
- Supported Terraform version: **v1.x**

This new Provider (ZTC Terraform Provider) enables fully automated provisioning and configuration across several key components in the Cloud & Branch Connector portal.

### NEW - DATA SOURCES

[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_activation_status`` - Activate ZTC Configuration information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_location_template`` - Retrieves location template information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_provisioning_url`` - Retrieves provisioning URL information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_location_management`` - Retrieves location management information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_edge_connector_group`` - Retrieves Edge Connector Group information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_traffic_forwarding_rule`` - Retrieves Traffic Forwarding rule
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_traffic_forwarding_dns_rule`` - Retrieves Traffic Forwarding DNS rule
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_traffic_forwarding_log_rule`` - Retrieves Traffic Forwarding LOG rule
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_forwarding_gateway`` - Retrieves ZIA and LOG Forwarding Gateway
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_dns_forwarding_gateway`` - Retrieves DNS Forwarding Gateway
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_ip_destination_groups`` - Retrieves IP Destination Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_ip_source_groups`` - Retrieves IP Source Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_ip_pool_groups`` - Retrieves IP Pool Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_network_services`` - Retrieves Network Services
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_network_service_groups`` - Retrieves Network Services Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_account_groups`` - Retrieves Account Groups
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_public_cloud_info`` - Retrieves Public Cloud Info
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_supported_regions`` - Retrieves Supported Regions
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Data Source ``ztc_workload_groups`` - Retrieves Workload Groups

### NEW - RESOURCE

[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_activation_status`` - Activate ZTC Configuration information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_location_template`` - Retrieves location template information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_provisioning_url`` - Retrieves provisioning URL information
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_traffic_forwarding_rule`` - Retrieves Traffic Forwarding rule
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_traffic_forwarding_dns_rule`` - Retrieves Traffic Forwarding DNS rule
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_traffic_forwarding_log_rule`` - Retrieves Traffic Forwarding LOG rule
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_forwarding_gateway`` - Retrieves ZIA and LOG Forwarding Gateway
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_dns_forwarding_gateway`` - Retrieves DNS Forwarding Gateway
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_ip_destination_groups`` - Retrieves IP Destination Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_ip_source_groups`` - Retrieves IP Source Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_ip_pool_groups`` - Retrieves IP Pool Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_network_services`` - Retrieves Network Services
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_network_service_groups`` - Retrieves Network Services Group
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_account_groups`` - Retrieves Account Groups
[PR #10](https://github.com/zscaler/terraform-provider-ztc/pull/4) - Resource ``ztc_public_cloud_info`` - Retrieves Public Cloud Info


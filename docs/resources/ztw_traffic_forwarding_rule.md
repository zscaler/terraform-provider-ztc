---
subcategory: "Policy Management"
layout: "zscaler"
page_title: "ZTW: traffic_forwarding_rule"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-forwarding-rules
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-rules
  Get information about Traffic Forwarding Rules.
---

# ztw_traffic_forwarding_rule (Resource)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-forwarding-rules)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-rules)

Use the **ztw_traffic_forwarding_rule** resource allows the creation and management of traffic forwarding in the Zscaler Cloud and Branch Connector Portal.

## Example Usage - Forward Method "DIRECT"

```hcl
resource "ztw_forwarding_gateway" "ztw_gw01" {
  name           = "ZTW_GW01"
  description    = "Example Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ZIA"
}

data "ztw_location_management" "this" {
  name = "AWS-CAN-ca-central-1-vpc-05c7f364cf47c2b93"
}

data "ztw_ip_destination_groups" "this" {
  name = "example"
}

data "zia_ip_source_groups" "this" {
  name = "example"
}

data "ztw_network_service" "this" {
  name = "ICMP_ANY"
}

data "ztw_network_services_groups" "this" {
  name = "Corporate Custom SSH TCP_10022"
}

# NOTE: To retrieve the src_workload_groups ID information, you must leverage the ZIA Terraform Provider at the moment,
# which returns the exact same information ID since the resource is cross-shared between ZIA and ZTW.

data "zia_workload_groups" "this" {
  name = "WORKLOAD_GROUP01"
}

resource "ztw_traffic_forwarding_rule" "this1" {
  name           = "DIRECT_Forwarding_Rule01"
  description    = "DIRECT Forwarding Rule 01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  type           = "EC_RDR"
  forward_method = "DIRECT"
  src_ips        = ["192.168.200.200"]
  dest_addresses = ["192.168.255.1"]
  wan_selection  = "BALANCED_RULE"
  dest_countries = ["CA", "US"]
  src_workload_groups {
    id = [data.zia_workload_groups.this.id]
  }
  nw_service_groups {
    id = [data.ztw_network_services_groups.this.id]
  }
  nw_services {
    id = [data.ztw_network_service.this.id]
  }
  src_ip_groups {
    id = [data.zia_ip_source_groups.this.id]
  }
  dest_ip_groups {
    id = [data.ztw_ip_destination_groups.this.id]
  }
  locations {
    id = [data.ztw_location_management.this.id]
  }
  proxy_gateway {
    id   = ztw_forwarding_gateway.ztw_gw01.id
    name = ztw_forwarding_gateway.ztw_gw01.name
  }
}
```

## Example Usage - Forward Method "ZIA"

```hcl
resource "ztw_forwarding_gateway" "ztw_gw01" {
  name           = "ZTW_GW01"
  description    = "Example Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ZIA"
}

data "ztw_location_management" "this" {
  name = "AWS-CAN-ca-central-1-vpc-05c7f364cf47c2b93"
}

data "ztw_ip_destination_groups" "this" {
  name = "example"
}

data "zia_ip_source_groups" "this" {
  name = "example"
}

data "ztw_network_service" "this" {
  name = "ICMP_ANY"
}

data "ztw_network_services_groups" "this" {
  name = "Corporate Custom SSH TCP_10022"
}

# NOTE: To retrieve the src_workload_groups ID information, you must leverage the ZIA Terraform Provider at the moment,
# which returns the exact same information ID since the resource is cross-shared between ZIA and ZTW.

data "zia_workload_groups" "this" {
  name = "WORKLOAD_GROUP01"
}

resource "ztw_traffic_forwarding_rule" "this1" {
  name           = "ZIA_Forwarding_Rule01"
  description    = "ZIA Forwarding Rule 01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  type           = "EC_RDR"
  forward_method = "ZIA"
  src_ips        = ["192.168.200.200"]
  dest_addresses = ["192.168.255.1"]
  wan_selection  = "BALANCED_RULE"
  dest_countries = ["CA", "US"]
  src_workload_groups {
    id = [data.zia_workload_groups.this.id]
  }
  nw_service_groups {
    id = [data.ztw_network_services_groups.this.id]
  }
  nw_services {
    id = [data.ztw_network_service.this.id]
  }
  src_ip_groups {
    id = [data.zia_ip_source_groups.this.id]
  }
  dest_ip_groups {
    id = [data.ztw_ip_destination_groups.this.id]
  }
  locations {
    id = [data.ztw_location_management.this.id]
  }
  proxy_gateway {
    id   = ztw_forwarding_gateway.ztw_gw01.id
    name = ztw_forwarding_gateway.ztw_gw01.name
  }
}
```

## Example Usage - Forward Method "ECZPA"

```hcl
resource "ztw_forwarding_gateway" "ztw_gw01" {
  name           = "ZTW_GW01"
  description    = "Example Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ZIA"
}

data "ztw_location_management" "this" {
  name = "AWS-CAN-ca-central-1-vpc-05c7f364cf47c2b93"
}

data "zia_ip_source_groups" "this" {
  name = "example"
}

# NOTE: To retrieve the src_workload_groups ID information, you must leverage the ZIA Terraform Provider at the moment,
# which returns the exact same information ID since the resource is cross-shared between ZIA and ZTW.

data "zia_workload_groups" "this" {
  name = "WORKLOAD_GROUP01"
}

resource "ztw_traffic_forwarding_rule" "this1" {
  name           = "ECZPA_Forwarding_Rule01"
  description    = "ECZPA Forwarding Rule 01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  type           = "EC_RDR"
  forward_method = "ECZPA"
  src_ips        = ["192.168.200.200"]
  wan_selection  = "BALANCED_RULE"
  src_workload_groups {
    id = [data.zia_workload_groups.this.id]
  }
  src_ip_groups {
    id = [data.zia_ip_source_groups.this.id]
  }
  locations {
    id = [data.ztw_location_management.this.id]
  }
  zpa_application_segments {
    id = [18612387, 18616051] # There is no current way to retrieve the ID information which is required for the configuration.
  }
  zpa_application_segment_groups {
    id = [18612386, 18661615] # There is no current way to retrieve the ID information which is required for the configuration.
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the forwarding rule to be exported.
* `id` - (Optional) The ID of the forwarding rule resource.
* `type` - (Optional) The rule type selected from the available options.
* `access_control` - (Optional) Access permission available for the current user to the rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) A unique identifier assigned to the forwarding rule.
* `name` - (String) The name of the forwarding rule.

### Optional

* `description` - (String) Additional information about the forwarding rule.
* `forward_method` - (String) The type of traffic forwarding method selected from the available options (e.g., DIRECT, ZIA, ECZPA, DROP, LOCAL_SWITCH).
* `state` - (String) Indicates whether the forwarding rule is enabled or disabled.
* `order` - (Number) The order of execution for the forwarding rule order.
* `rank` - (Number) Admin rank assigned to the forwarding rule.
* `type` - (String) The rule type (e.g., FIREWALL, DNS, DNAT, SNAT, FORWARDING, INTRUSION_PREVENTION, EC_DNS, EC_RDR, EC_SELF, DNS_RESPONSE).
* `src_ips` - (List of String) User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
* `dest_addresses` - (List of String) List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses.
* `dest_countries` - (List of String) Destination countries for which the rule is applicable.
* `wan_selection` - (String) WAN selection (only applicable when configuring a hardware device deployed in gateway mode).
* `source_ip_group_exclusion` - (Boolean) Source IP groups that must be excluded from the rule application.
* `ec_groups` - (List of Object) Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `locations` - (List of Object) Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `src_ip_groups` - (List of Object) Source IP address groups for which the rule is applicable.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `dest_ip_groups` - (List of Object) User-defined destination IP address groups to which the rule is applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `nw_services` - (List of Object) User-defined network services to which the rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `nw_service_groups` - (List of Object) User-defined network service groups to which the rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `app_service_groups` - (List of Object) List of application service groups.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `zpa_application_segments` - (List of Object) List of ZPA Application Segments for which this rule is applicable (used for ECZPA forwarding method).
  * `id` - (Number) Application segment identifier.
* `zpa_application_segment_groups` - (List of Object) List of ZPA Application Segment Groups for which this rule is applicable (used for ECZPA forwarding method).
  * `id` - (Number) Application segment group identifier.
* `proxy_gateway` - (List of Object) The proxy gateway for which the rule is applicable (for Proxy Chaining forwarding method).
  * `id` - (Number) Gateway identifier.
  * `name` - (String) Gateway name.
* `src_workload_groups` - (List of Object) The list of preconfigured workload groups to which the policy must be applied.
  * `id` - (Number) Workload group identifier.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTW configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztw_traffic_forwarding_rule** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import ztw_traffic_forwarding_rule.example <rule_id>
```

or

```shell
terraform import ztw_traffic_forwarding_rule.example <rule_name>
```

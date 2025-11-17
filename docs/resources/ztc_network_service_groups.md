---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTC: network_service_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-network-service-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/network-service-groups
  Creates and manages Network Service Groups.
---

# ztc_network_service_groups (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/policy-resources#/networkServiceGroups-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-network-service-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/network-service-groups)

Use the **ztc_network_service_groups** resource allows the creation and management of network service groups available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTC traffic forwarding rule.

## Example Usage - Create Network Services Group

```hcl
data "ztc_network_service" "example" {
  name = "ICMP_ANY"
}


resource "ztc_network_services_groups" "example" {
  name        = "example"
  description = "example"
  services {
    id = [
      data.ztc_network_service.example.id,
    ]
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the network service group to be exported.
* `id` - (Optional) The ID of the network service group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) ID of the network service group.
* `name` - (String) Name of the network service group.

### Optional

* `description` - (String) Description of the network service group.
* `services` - (List of Object) List of network services in the group.
  * `id` - (Number) ID of network service.
  * `name` - (String) Name of network service.
  * `tag` - (String) Tag associated with the network service.
  * `type` - (String) Type of network service: standard, predefined, or custom.
  * `description` - (String) Description of network service.
  * `is_name_l10n_tag` - (Number) Indicates the external ID. Applicable only when this reference is of an external entity.
  * `src_tcp_ports` - (List of Object) Source TCP ports.
    * `start` - (Number) Starting port number (1-65535).
    * `end` - (Number) Ending port number (1-65535).
  * `dest_tcp_ports` - (List of Object) Destination TCP ports.
    * `start` - (Number) Starting port number (1-65535).
    * `end` - (Number) Ending port number (1-65535).
  * `src_udp_ports` - (List of Object) Source UDP ports.
    * `start` - (Number) Starting port number (1-65535).
    * `end` - (Number) Ending port number (1-65535).
  * `dest_udp_ports` - (List of Object) Destination UDP ports.
    * `start` - (Number) Starting port number (1-65535).
    * `end` - (Number) Ending port number (1-65535).

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_network_service_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import ztc_network_service_groups.example <rule_id>
```

or

```shell
terraform import ztc_network_service_groups.example <rule_name>
```

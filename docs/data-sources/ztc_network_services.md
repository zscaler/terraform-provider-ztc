---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTC: network_services"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-network-services
  API documentation https://help.zscaler.com/cloud-branch-connector/network-services
  Get information about Network Services.
---

# ztc_network_services (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-network-services)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/network-services)

Use the **ztc_network_services** data source to get information about network services available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with with ZTC traffic forwarding rule.

## Example Usage - Retrieve by Name

```hcl
data "ztc_network_services" "example" {
    name = "example_network_service"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_network_services" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the network service to be exported.
* `id` - (Optional) The ID of the network service resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) ID of network service.
* `name` - (String) Name of network service.

### Optional

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
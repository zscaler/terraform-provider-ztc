---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTC: network_services"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-network-services
  API documentation https://help.zscaler.com/cloud-branch-connector/network-services
  Creates and manages Network Services.
---

# ztc_network_services (Resource)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-network-services)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/network-services)

Use the **ztc_network_services** resource allows the creation and management of network services available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTC traffic forwarding rule.

## Example Usage - Create Network Services

```hcl
resource "ztc_network_services" "example" {
  name        = "example"
  description = "example"
  src_tcp_ports {
    start = 5000
  }
  src_tcp_ports {
    start = 5001
  }
  src_tcp_ports {
    start = 5002
    end   = 5005
  }
  dest_tcp_ports {
    start = 5000
  }
  dest_tcp_ports {
    start = 5001
  }
  dest_tcp_ports {
    start = 5003
    end   = 5005
  }
  type = "CUSTOM"
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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_network_services** can be imported by using `<SERVICE_ID>` or `<SERVICE_NAME>` as the import ID.

For example:

```shell
terraform import ztc_network_services.example <rule_id>
```

or

```shell
terraform import ztc_network_services.example <rule_name>
```

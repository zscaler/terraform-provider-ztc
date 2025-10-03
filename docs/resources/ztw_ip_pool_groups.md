---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTW: ip_pool_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-ip-pool-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/ip-pool-groups
  Creates and manages IP Pool Groups.
---

# ztw_ip_pool_groups (Resource)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-ip-pool-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/ip-pool-groups)

Use the **ztw_ip_pool_groups** resource allows the creation and management of IP Pool Groups available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTW traffic forwarding rule.

## Example Usage - Create IP Pool Group

```hcl
resource "ztw_ip_pool_groups" "example" {
  name        = "Example IP Pool"
  description = "Updated Example IP Pool for testing"
  ip_addresses = [
    "192.168.100.0/24"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the IP pool group to be exported.
* `id` - (Optional) The ID of the IP pool group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) ID of the IP address group or IP pool.
* `name` - (String) Name of the IP address group or IP pool.

### Optional

* `description` - (String) Description of the IP group or IP pool.
* `ip_addresses` - (List of String) IP Subnets included in the IP group or IP pool. Only `ONE` CIDR subnet is allowed i.e `10.0.0.0/24`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTW configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztw_ip_pool_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import ztw_ip_pool_groups.example <rule_id>
```

or

```shell
terraform import ztw_ip_pool_groups.example <rule_name>
```
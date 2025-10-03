---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTW: ip_source_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-ip-source-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/ip-source-groups
  Creates and manages IP Source Groups.
---

# ztw_ip_source_groups (Resource)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-ip-source-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/ip-source-groups)

Use the **ztw_ip_source_groups** resource allows the creation and management of IP Source Groups available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTW traffic forwarding rule.

## Example Usage - Create IP Source Group

```hcl
resource "zia_ip_source_groups" "this" {
  name         = "example1"
  description  = "example1"
  ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the IP source group to be exported.
* `id` - (Optional) The ID of the IP source group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) ID of the IP address group or IP pool.
* `name` - (String) Name of the IP address group or IP pool.

### Optional

* `description` - (String) Description of the IP group or IP pool.
* `ip_addresses` - (List of String) IP addresses included in the IP group or IP pool.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTW configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztw_ip_source_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import ztw_ip_source_groups.example <rule_id>
```

or

```shell
terraform import ztw_ip_source_groups.example <rule_name>
```
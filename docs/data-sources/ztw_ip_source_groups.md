---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTW: ip_source_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-ip-source-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/ip-source-groups
  Get information about IP Source Groups.
---

# ztw_ip_source_groups (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-ip-source-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/ip-source-groups)

Use the **ztw_ip_source_groups** data source to get information about IP Source Groups available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with ZTW traffic forwarding rule.

## Example Usage - Retrieve by Name

```hcl
data "ztw_ip_source_groups" "example" {
    name = "example_ip_source_group"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztw_ip_source_groups" "example" {
    id = 5458452
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
* `creator_context` - (String) Indicates that the IP group or IP pool is created in Cloud & Branch Connector (EC) (only applicable value).
* `is_non_editable` - (Boolean) Indicates whether the group is view-only (true) or editable (false).
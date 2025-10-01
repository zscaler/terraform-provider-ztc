---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTW: ip_destination_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-ip-destination-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/ip-destination-groups
  Get information about IP Destination Groups.
---

# ztw_ip_destination_groups (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-ip-destination-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/ip-destination-groups)

Use the **ztw_ip_destination_groups** data source to get information about IP Destination Groups available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with ZTW traffic forwarding rule.

## Example Usage - Retrieve by Name

```hcl
data "ztw_ip_destination_groups" "example" {
    name = "example_ip_destination_group"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztw_ip_destination_groups" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the IP destination group to be exported.
* `id` - (Optional) The ID of the IP destination group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) ID of the destination IP group.
* `name` - (String) Name of the destination IP group.

### Optional

* `description` - (String) Description of the group.
* `type` - (String) Type of the destination IP group (e.g., DSTN_IP, DSTN_FQDN, DSTN_DOMAIN, DSTN_OTHER).
* `ip_addresses` - (List of String) IP addresses or domain names included in the group.
* `countries` - (List of String) Countries included in the group.
* `is_non_editable` - (Boolean) Indicates whether the group is view-only (true) or editable (false).
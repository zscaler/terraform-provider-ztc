---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTW: ip_destination_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-ip-destination-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/ip-destination-groups
  Creates and manages IP Destination Groups.
---

# ztw_ip_destination_groups (Resource)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-ip-destination-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/ip-destination-groups)

Use the **ztw_ip_destination_groups** resource allows the creation and management of IP Destination Groups available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTW traffic forwarding rule.

## Example Usage - Type DSTN_FQDN

```hcl
resource "ztw_ip_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses   = ["test1.acme.com", "test2.acme.com", "test3.acme.com"]
}
```

## Example Usage - Type DSTN_IP

```hcl
resource "ztw_ip_destination_groups" "example_ip_ranges" {
  name        = "Example - IP Ranges"
  description = "Example - IP Ranges"
  type        = "DSTN_IP"
  addresses = ["3.217.228.0-3.217.231.255",
    "3.235.112.0-3.235.119.255",
    "52.23.61.0-52.23.62.25",
  "35.80.88.0-35.80.95.255"]
}
```

## Example Usage - Type DSTN_DOMAIN

```hcl
resource "ztw_ip_destination_groups" "example_dstn_domain" {
  name        = "Example Destination Domain"
  description = "Example Destination Domain"
  type        = "DSTN_DOMAIN"
  addresses   = ["acme.com", "acme1.com"]
}
```

## Example Usage - Type DSTN_OTHER

```hcl
resource "ztw_ip_destination_groups" "example_dstn_other" {
  name        = "Example Destination Other"
  description = "Example Destination Other"
  type        = "DSTN_OTHER"
  countries   = ["CA"]
}
```

## Example Usage - 
```hcl
resource "ztw_ip_destination_groups" "example" {
  name        = "Example"
  description = "Example"
  type        = "DSTN_FQDN"
  addresses   = ["test1.acme.com", "test2.acme.com", "test3.acme.com"]
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
* `type` - (String) Type of the destination IP group (e.g., `DSTN_IP`, `DSTN_FQDN`, `DSTN_DOMAIN`, DSTN_OTHER).
* `ip_addresses` - (List of String) IP addresses or domain names included in the group.
* `countries` - (List of String) The list of countries that must be included in the rule based on the rule. If no value is set, this field is ignored during policy evaluation and the rule is applied to all source countries.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTW configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztw_ip_destination_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import ztw_ip_destination_groups.example <rule_id>
```

or

```shell
terraform import ztw_ip_destination_groups.example <rule_name>
```
---
subcategory: "Partner Integrations"
layout: "zscaler"
page_title: "ZTC: account_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-amazon-web-services-account-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/partner-integrations#/accountGroups-get
  Retrieves the details of AWS account groups with metadata
---

# ztc_account_groups (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/partner-integrations#/accountGroups-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-amazon-web-services-account-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/partner-integrations#/accountGroups-get)

Use the **ztc_account_groups** data source to get information details of AWS account groups with metadata.

## Example Usage - Retrieve by Name

```hcl
data "ztc_account_groups" "example" {
    name = "AWS_Account_Group01"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_account_groups" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Number) The ID of the AWS account group.
* `name` - (String) The name of the AWS account group. Must be non-null, non-empty, unique, and 128 characters or fewer in length.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Optional

* `description` - (String) The description of the AWS account group. Must be less than or equal to 512 characters.
* `cloud_type` - (String) The cloud type. The default and manadatory value is AWS. Returned values are: `AWS`, `AZURE`, `GCP`
* `cloud_connector_groups` - (List of Object)
  * `id` - (Number) An ID that uniquely identifies an entity.
* `public_cloud_accounts` - (List of Object)
  * `id` - (Number) An ID that uniquely identifies an entity.

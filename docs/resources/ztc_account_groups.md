---
subcategory: "Partner Integrations"
layout: "zscaler"
page_title: "ZTC: account_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-amazon-web-services-account-groups
  API documentation https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zcloudconnector/partner-integrations/aws-account-group-z-resource-create-account-group
  Creates an AWS account group. You can create a maximum of 128 groups in each organization.
---

# ztc_account_groups (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/partner-integrations#/accountGroups-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-amazon-web-services-account-groups)
* [API documentation](https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zcloudconnector/partner-integrations/aws-account-group-z-resource-create-account-group)

The **ztc_account_groups** resource allows you to create and manage Account Group configurations in the Zscaler Zero Trust Cloud (ZTC) platform.

## Example Usage - Create a Account Group without Cloud Connector Group

```hcl
resource "ztc_account_groups" "this" {
    name        = "AWS_Account_Group01"
    description = "AWS_Account_Group01"
    cloud_type  = "AWS"
    public_cloud_accounts {
      id = [2815549]
    }
}
```

## Example Usage - Create a Account Group with Cloud Connector Group

```hcl
resource "ztc_account_groups" "this" {
    name        = "AWS_Account_Group01"
    description = "AWS_Account_Group01"
    cloud_type  = "AWS"
    public_cloud_accounts {
      id = [2815549]
    }
    cloud_connector_groups {
      id = [2815559]
    }
}
```

## Argument Reference

### Required
* `name` - (String) The name of the AWS account group. Must be non-null, non-empty, unique, and 128 characters or fewer in length.

### Optional

* `description` - (String) The description of the AWS account group. Must be less than or equal to 512 characters.
* `cloud_type` - (String) The cloud type. The default and manadatory value is AWS. Returned values are: `AWS`, `AZURE`, `GCP`
* `cloud_connector_groups` - (List of Object)
  * `id` - (Number) An ID that uniquely identifies an entity.
* `public_cloud_accounts` - (List of Object)
  * `id` - (Number) An ID that uniquely identifies an entity.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_account_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import ztc_account_groups.example <gateway_id>
```

or

```shell
terraform import ztc_account_groups.example <gateway_name>
```
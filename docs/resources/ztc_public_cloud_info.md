---
subcategory: "Partner Integrations"
layout: "zscaler"
page_title: "ZTC: public_cloud_info"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/adding-amazon-web-services-account
  API documentation https://help.zscaler.com/cloud-branch-connector/partner-integrations#/publicCloudInfo-post
  Creates a new AWS account with the provided account and region details
---

# ztc_public_cloud_info (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/partner-integrations#/publicCloudInfo-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/adding-amazon-web-services-account)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/partner-integrations#/publicCloudInfo-post)

Use the **ztc_public_cloud_info** resource allows the creates a new AWS account with the provided account and region details in the Zscaler Cloud and Branch Connector Portal.

## Example Usage - Public Cloud Info

```hcl
data "ztc_supported_regions" "this" {
  name = "US_EAST_1"
}

resource "ztc_public_cloud_info" "this" {
  name       = "AWSAccount01"
  cloud_type = "AWS"
  account_details {
    aws_account_id           = "123456789"
    aws_role_name            = "zscaler-role"
    cloud_watch_group_arn    = "DISABLED"
    event_bus_name           = "zscaler-bus-123456-zscalerthree.net"
    trouble_shooting_logging = true
    trusted_account_id       = "123456789"
    trusted_role             = "arn:aws:iam::123456789:role/ZscalerTagDiscoveryRole"
  }
  supported_regions {
    id = [data.ztc_supported_regions.this.id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the IP pool group to be exported.
* `id` - (Optional) The ID of the IP pool group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

## Attribute Reference

* `cloud_type` - (String) The cloud type. The default and mandatory value is AWS. Supported values are: `AWS`, `AZURE`, `GCP`

### account_groups

* `account_groups` - (List of Object) An immutable reference to account groups, which consists of ID and name.
  * `id` - (Number) An ID that uniquely identifies an entity.

### supported_regions

* `supported_regions` - (List of Object) Regions supported by Zscaler's Tag Discovery Service.
  * `id` - (Number) The unique ID of the supported region.

### account_details

* `account_details` - (List of Object) The AWS account details.
  * `name` - (String) The name of the AWS account.
  * `aws_account_id` - (String) The AWS account ID where workloads are deployed. The ID is non-null, non-empty, and unique, and contains 12 digits.
  * `aws_role_name` - (String) The AWS trusting role in your account. The name is non-null, non-empty, and 64 characters or fewer in length.
  * `cloud_watch_group_arn` - (String) The resource name (ARN) of the AWS CloudWatch log group.
  * `event_bus_name` - (String) The name of the event bus that sends notifications to the Zscaler service using EventBridge.
  * `external_id` - (String) The unique external ID for the AWS account. If provided, it must match the externalId specified outside of accountDetails.
  * `log_info_type` - (String) The type of log information. Supported types are `INFO` and `ERROR`.
  * `trouble_shooting_logging` - (Boolean) Indicates whether logging is enabled for troubleshooting purposes.
  * `trusted_account_id` - (String) The ID of the Zscaler AWS account.
  * `trusted_role` - (String) The name of the trusted role in the Zscaler AWS account.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_public_cloud_info** can be imported by using `<CLOUD_ID>` or `<CLOUD_NAME>` as the import ID.

For example:

```shell
terraform import ztc_public_cloud_info.example <rule_id>
```

or

```shell
terraform import ztc_public_cloud_info.example <rule_name>
```
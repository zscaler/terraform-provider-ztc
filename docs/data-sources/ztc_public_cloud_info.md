---
subcategory: "Partner Integrations"
layout: "zscaler"
page_title: "ZTC: public_cloud_info"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/adding-amazon-web-services-account
  API documentation https://help.zscaler.com/cloud-branch-connector/partner-integrations#/publicCloudInfo-get
  Retrieves the list of AWS accounts with metadata.
---

# ztc_public_cloud_info (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/adding-amazon-web-services-account)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/partner-integrations#/publicCloudInfo-get)

Use the **ztc_public_cloud_info** data source to get information details of AWS accounts with metadata, including account details, supported regions, region status, and associated account groups.

## Example Usage - Retrieve by Name

```hcl
data "ztc_public_cloud_info" "example" {
    name = "AWSAccountInfo01"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_public_cloud_info" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The unique ID of the AWS account.
* `name` - (Optional) The name of the AWS account. Must be non-null, non-empty, unique, and 128 characters or fewer in length.

-> **NOTE:** At least one of `id` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cloud_type` - (String) The cloud type. The default and mandatory value is AWS. Supported values are: `AWS`, `AZURE`, `GCP`
* `external_id` - (String) A unique external ID for the AWS account.
* `last_mod_time` - (Number) The date and time when the AWS account was last modified (Unix timestamp).
* `last_sync_time` - (Number) The last time the AWS account was synced (Unix timestamp).

### account_groups

* `account_groups` - (List of Object) An immutable reference to account groups, which consists of ID and name.
  * `id` - (Number) An ID that uniquely identifies an entity.
  * `name` - (String) The name of the entity.

### last_mod_user

* `last_mod_user` - (List of Object) Automatically populated with the current ZTC admin user, after a successful POST or PUT request.
  * `id` - (Number) An ID that uniquely identifies an entity.
  * `name` - (String) The name of the entity.

### region_status

* `region_status` - (List of Object) The status and configuration details of the region where the workloads are deployed.
  * `id` - (Number) The unique ID of the region.
  * `name` - (String) The name of the region.
  * `cloud_type` - (String) The cloud type. The default and mandatory value is AWS. Supported values: `AWS`, `AZURE`, `GCP`
  * `status` - (Boolean) Indicates the operational status of the region.

### supported_regions

* `supported_regions` - (List of Object) Regions supported by Zscaler's Tag Discovery Service.
  * `id` - (Number) The unique ID of the supported region.
  * `name` - (String) The name of the supported region.
  * `cloud_type` - (String) The cloud type. The default and mandatory value is AWS. Supported values: `AWS`, `AZURE`, `GCP`

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

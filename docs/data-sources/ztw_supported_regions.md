---
subcategory: "Partner Integrations"
layout: "zscaler"
page_title: "ZTW: supported_regions"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/adding-amazon-web-services-account
  API documentation https://help.zscaler.com/cloud-branch-connector/partner-integrations#/publicCloudInfo/supportedRegions-getpublicCloudInfo-get
  Retrieves a list of AWS regions supported for workload discovery settings (WDS).
---

# ztw_supported_regions (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/adding-amazon-web-services-account)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/partner-integrations#/publicCloudInfo/supportedRegions-get)

Use the **ztw_supported_regions** data source to get information details AWS regions supported for workload discovery settings (WDS). This data source can be used when configuring the attribute `supported_regions` in the resource `ztw_public_cloud_info`

## Example Usage - Retrieve All Regions

```hcl
data "ztw_supported_regions" "all" {
}

output "all_regions" {
  value = data.ztw_supported_regions.all.regions
}
```

## Example Usage - Retrieve by Name

```hcl
data "ztw_supported_regions" "example" {
    name = "US_EAST_1"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztw_supported_regions" "example" {
    id = 1178341
}
```

## Argument Reference

All arguments are optional:

* `id` - (Optional) The unique ID of the supported region. When specified, returns a single region.
* `name` - (Optional) The name of the supported region. When specified, returns a single region.

-> **NOTE:** When neither `id` nor `name` is specified, all supported regions are returned in the `regions` attribute. When either `id` or `name` is specified, only that single region's details are returned in the top-level attributes.

## Attribute Reference

The following attributes are exported:

### Single Region Attributes

These attributes are populated when `id` or `name` is specified:

* `cloud_type` - (String) The cloud type. The default value is AWS. Supported values are: `AWS`, `AZURE`, `GCP`

### All Regions Attribute

This attribute is populated when neither `id` nor `name` is specified:

* `regions` - (List of Object) List of all supported regions.
  * `id` - (Number) The unique ID of the supported region.
  * `name` - (String) The name of the supported region.
  * `cloud_type` - (String) The cloud type. Supported values are: `AWS`, `AZURE`, `GCP`

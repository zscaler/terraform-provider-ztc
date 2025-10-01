# Location Template Example

This example will show you how to use Terraform to create and retriieve location template information in the Zscaler Cloud and Branch Connector PortalPortal
This example codifies [this API](https://help.zscaler.com/cloud-branch-connector/about-ip-source-groups).

To run, configure your ZTW provider as described [Here](https://github.com/zscaler/terraform-provider-ztw/blob/master/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```

## Destroy 💥

```bash
terraform destroy
```

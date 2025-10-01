# Network Service Groups Example

This example will show you how to use Terraform to create a network service groups available in the Zscaler Cloud and Branch Connector Portal
This example codifies [this API](https://help.zscaler.com/cloud-branch-connector/about-network-service-groups).

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

# Location Management

This example will show you how to use Terraform to retrieve Location Management available in the Zscaler Cloud and Branch Connector Portal
This example codifies [this API](https://help.zscaler.com/cloud-branch-connector/about-locations).

To run, configure your ZTC provider as described [Here](https://github.com/zscaler/terraform-provider-ztc/blob/master/docs/index.html.markdown)

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

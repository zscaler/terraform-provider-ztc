---
page_title: "Troubleshooting Guide"
---

# How to troubleshoot your problem

If you have problems with code that uses ZPA Terraform provider, follow these steps to solve them:

* Check symptoms and solutions in the [Typical problems](#typical-problems) section below.
* Upgrade provider to the latest version. The bug might have already been fixed.
* In case of authentication problems, see the [Authentication Issues](#authentication-issues) below.
* Collect debug information using following command:

```sh
TF_LOG=DEBUG ZSCALER_SDK_VERBOSE=true ZSCALER_SDK_LOG=true terraform apply -no-color 2>&1 |tee tf-debug.log
```

* Open a [new GitHub issue](https://github.com/zscaler/terraform-provider-ztc/issues/new/choose) providing all information described in the issue template - debug logs, your Terraform code, Terraform & plugin versions, etc.

## Typical problems

### Authentication Issues

### │ Error: Invalid provider configuration and Error: failed configuring the provider

The most common problem with invalid provider is when the ZTC API credentials are not properly set via one of the supported methods. Please make sure to read the documentation for the supported authentication methods [Authentication Methods](https://registry.terraform.io/providers/zscaler/ztc/latest/docs)

```sh
│ Provider "zscaler/ztc" requires explicit configuration. Add a provider block to the root module and configure the
│ provider's required arguments as described in the provider documentation.
```

```sh
│ Error: expected ztc_cloud to be one of ["zscaler" "zscalerone" "zscalertwo" "zscalerthree" "zscloud" "zscalerbeta" "zscalergov" "zscalerten" "zspreview"], got
│
│   with provider["zscaler.com/ztc/ztc"],
│   on <input-prompt> line 1:
│   (source code not available)
│
```

## Multiple Provider Configurations

The most common reason for technical difficulties might be related to missing `alias` attribute in `provider "ztc" {}` blocks or `provider` attribute in `resource "ztc_..." {}` blocks, when using multiple provider configurations. Please make sure to read [`alias`: Multiple Provider Configurations](https://www.terraform.io/docs/language/providers/configuration.html#alias-multiple-provider-configurations) documentation article.

## Error while installing: registry does not have a provider

```sh
Error while installing hashicorp/ztc: provider registry
registry.terraform.io does not have a provider named
registry.terraform.io/hashicorp/ztc
```

If you notice below error, it might be due to the fact that [required_providers](https://www.terraform.io/docs/language/providers/requirements.html#requiring-providers) block is not defined in *every module*, that uses ZTC Terraform Provider. Create `versions.tf` file with the following contents:

```hcl
# versions.tf
terraform {
  required_providers {
    ztc = {
      source  = "zscaler/ztc"
      version = "0.1.0"
    }
  }
}
```

... and copy the file in every module in your codebase. Our recommendation is to skip the `version` field for `versions.tf` file on module level, and keep it only on the environment level.

```
├── environments
│   ├── sandbox
│   │   ├── README.md
│   │   ├── main.tf
│   │   └── versions.tf
│   └── production
│       ├── README.md
│       ├── main.tf
│       └── versions.tf
└── modules
    ├── first-module
    │   ├── ...
    │   └── versions.tf
    └── second-module
        ├── ...
        └── versions.tf
```

## Error: Failed to install provider

Running the `terraform init` command, you may see `Failed to install provider` error if you didn't check-in [`.terraform.lock.hcl`](https://www.terraform.io/language/files/dependency-lock#lock-file-location) to the source code version control:

```sh
Error: Failed to install provider

Error while installing zscaler/ztc: v0.1.0: checksum list has no SHA-256 hash for "https://github.com/zscaler/terraform-provider-ztc/releases/download/v1.0.0/terraform-provider-ztc.0.1.0_darwin_amd64.zip"
```

You can fix it by following three simple steps:

* Replace `zscaler.com/ztc/ztc` with `zscaler/ztc` in all your `.tf` files with the `python3 -c "$(curl -Ls https://github.com/zscaler/terraform-provider-ztc/scripts/upgrade-namespace.py)"` command.
* Run the `terraform state replace-provider zscaler.com/ztc/ztc zscaler/ztc` command and approve the changes. See [Terraform CLI](https://www.terraform.io/cli/commands/state/replace-provider) docs for more information.
* Run `terraform init` to verify everything working.

The terraform apply command should work as expected now.

## Error: Failed to query available provider packages

See the same steps as in [Error: Failed to install provider](#error-failed-to-install-provider).

### Error: Provider registry.terraform.io/zscaler/ztc v... does not have a package available for your current platform, windows_386

This kind of errors happens when the 32-bit version of ZTC Terraform provider is used, usually on Microsoft Windows. To fix the issue you need to switch to use of the 64-bit versions of Terraform and ZTC Terraform provider.

## Error: Unauthorized access to private resource

This error indicates that the resource is NOT accessible via the Legacy API and only via OneAPI. See documentation notes in each resource.

```sh
│ Error: error retrieving forwarding control rule 1263518: Error: {
│   "code": null,
│   "message": "Unauthorized access to private resource",
│   "url": "https://connector.zscalerbeta.net/api/v1/ecRules/ecRdr/1263518",
│   "status": 403
│ }
```
# Terraform Provider GitHub v4 (GraphQL)

For use alongside `terraform-provider-github` while it is determined how v4 resources will be included.

## Installation

Copy the built binaries into the Terraform working directory:

`build/darwin/terraform-provider-github-v4_vx.x.x` -> `/terraform.d/plugins/darwin_amd64/terraform-provider-github-v4_vx.x.x`

`build/linux/terraform-provider-github-v4_vx.x.x` -> `/terraform.d/plugins/linux_amd64/terraform-provider-github-v4_vx.x.x`

Read more about [plugin locations](https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations).

## Usage

```hcl
# The official Github provider
provider "github" {
  organization = "myorg"
  ...
}

# The v4 Github provider
provider "github-v4" {
  organization = "myorg"
  ...
}

# For data sources and resources you wish to use the v4 provider
resource "github_branch_protection" "master" {
  provider = github-v4
  ...
}
```

There are schema differences between the providers. For now you'll need to view the source.
 
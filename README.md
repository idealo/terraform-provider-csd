![Maintained](https://img.shields.io/maintenance/yes/2025) [![Release](https://github.com/idealo/terraform-provider-csd/actions/workflows/release.yml/badge.svg)](https://github.com/idealo/terraform-provider-csd/actions/workflows/release.yml) [![Test](https://github.com/idealo/terraform-provider-csd/actions/workflows/test.yml/badge.svg)](https://github.com/idealo/terraform-provider-csd/actions/workflows/test.yml)

# Terraform provider for [Common Short Domain product](https://github.com/idealo/transport_csd)

The Common Short Domain product gives you cool short domains (AWS Hosted Zones) in your AWS account so you can manage them yourself, without the hassle of a third party.

_Keep in mind that your FQDN shouldn't exceed 64 characters (including the final dot) to retrieve a TLS certificate._

# Installation

## Install from Terraform Registry

You can find our Terraform provider in the [Terraform registry](https://registry.terraform.io/providers/idealo/csd/latest).

Online documentation can also be found [here](https://registry.terraform.io/providers/idealo/csd/latest/docs).

## Upgrade from v1.x to v2.x

1. Comment all old `csd_zone` resources
2. Run `terraform apply`, this will delete your old zone delegation
3. Update provider version to `~>2.0`
4. Uncomment and rename old `csd_zone` resources to `csd_zone_delegation`
5. Run `terraform init --upgrade` to install the new version
6. Run `terraform apply` to put the zone delegations back in place

**❗Note: Your zone delegation will not work between steps 2 and 6. The DNS systems caches should cover this short downtime.**

# Usage

```terraform
terraform {
  required_version = "~>1.3"
  required_providers {
    csd = {
      source  = "idealo/csd"
      version = "~>2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~>4.8"
    }
  }
  backend "s3" {
    bucket         = "<ENTER_BUCKET_NAME>"
    key            = "global/s3/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}

# Setup AWS provider
provider "aws" {
  region              = "eu-central-1"
  allowed_account_ids = ["<ENTER_ACCOUNT_ID>"]
}

# Setup csd provider
# It will use the AWS credentials provided by environment variables or parameters
# The OIDC provider sets up the neccessary environment variables by default
provider "csd" {}

# Setup OIDC provider
# https://confluence.idealo.cloud/pages/viewpage.action?spaceKey=PTN&title=How+to+authenticate+from+GitHub+to+AWS
module "terraform_execution_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role"
  version = "~> 6.0"

  name            = "<ENTER_ROLE_NAME>"
  use_name_prefix = false
  max_session_duration = 6 * 60 * 60

  enable_github_oidc = true
  oidc_audiences     = ["sts.amazonaws.com"]
  oidc_wildcard_subjects = [
    "repo:idealo/<ENTER_REPO_NAME>:*",
  ]

  policies = {
    <ENTER_POLICY_NAME> = "arn:aws:iam::aws:policy/<ENTER_POLICY_NAME>"
  }
}
```

## Hosted zone delegation

```terraform
# Create a Route53 Hosted Zone.
# sample-app is a placeholder for the subdomain for your application.
# example.net is a placeholder for a domain which is supported in the CSD product.
resource "aws_route53_zone" "sample-app" {
  name = "sample-app.example.net"
}

# Create zone delegation in example.net zone via CSD provider
# example.net is a placeholder for a domain which is supported in the CSD product.
resource "csd_zone_delegation" "sample-app" {
  name         = aws_route53_zone.sample-app.name
  name_servers = aws_route53_zone.sample-app.name_servers
}
```

**⚠️ Important:** Keep in mind that the TTL of the NS records for your Hosted Zone can be up to 2 days. So destroying them could lead to extended downtimes for your workloads. We suggest to protect them as shown in the example above and/or separate their automation completely from your product workloads.

# FAQ

## Q: Provider does not support resource type

If you see the following error message, you accidentally updated from version 1.x to version 2.x:

```
│ Error: Invalid resource type
│
│   on main.tf line 42, in resource "csd_zone" "my_zone_delegation":
│   27: resource "csd_zone" "my_zone_delegation" {}
│
│ The provider idealo/csd does not support resource type "csd_zone".
```

To fix this issue downgrade to version 1.x like this:

```terraform
terraform {
  required_providers {
    csd = {
      source  = "idealo/csd"
      version = "~>1.0"
    }
  }
}
```

Now follow the proper upgrade procedure described [here](https://github.com/idealo/terraform-provider-csd/tree/main#upgrade-from-v1x-to-v2x).

# Development

For development notes see [DEV.md](DEV.md).




---

Made with ❤️ and ✨ by [🌐 Team Transport](https://github.com/orgs/idealo/teams/transport).

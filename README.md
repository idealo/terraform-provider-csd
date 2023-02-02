![Maintained](https://img.shields.io/maintenance/yes/2023)

# Terraform provider for [Common Short Domain product](https://github.com/idealo/transport_csd)

The Common Short Domain product gives you cool short domains (AWS Hosted Zones) in your AWS account so you can manage them yourself, without the hassle of a third party.

Currently we support the following domains where you can get subdomains:

- `idealo.tools`: Used for internal idealo tooling

More domains will follow in future updates. If you're missing one that you need, contact Team Transport.

_Keep in mind that your FQDN shouldn't exceed 64 characters (including the final dot) to retrieve a TLS certificate._

# Installation

## Install from Terraform Registry

_TBA_

## Manual Installation

1. Clone git repository
   ```shell
   git clone git@github.com:idealo/terraform-provider-csd.git
   ```
2. Build and install
   ```shell
   make install
   ```

# Usage

```terraform
terraform {
  required_version = "~> 1.3"
  required_providers {
    csd = {
      source  = "idealo.com/transport/csd"
      version = "~>1.0"
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
# https://confluence.eu.idealo.com/pages/viewpage.action?spaceKey=PTN&title=How+to+authenticate+from+GitHub+to+AWS
module "terraform_execution_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "~> 4.3"

  create_role = true
  role_name = "<ENTER_ROLE_NAME>"
  max_session_duration = 6 * 60 * 60

  provider_url = "token.actions.githubusercontent.com"
  oidc_subjects_with_wildcards = [
    "repo:idealo/<ENTER_REPO_NAME>:*",
  ]

  role_policy_arns = [
    "arn:aws:iam::aws:policy/<ENTER_POLICY_NAME>",
  ]
  number_of_role_policy_arns = 1
}

# Create desired zone in Route53
resource "aws_route53_zone" "shopverwaltung" {
  name = "shopverwaltung.idealo.tools"
}

# Create zone forwarding in idealo.tools zone via CSD provider
resource "csd_zone" "shopverwaltung" {
  name         = aws_route53_zone.shopverwaltung.name
  name_servers = aws_route53_zone.shopverwaltung.name_servers
}
```

# Development

For development notes see [DEV.md](DEV.md).

---

Made with ❤️ and ✨ by [🌐 Team Transport](https://github.com/orgs/idealo/teams/transport).


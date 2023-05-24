![Maintained](https://img.shields.io/maintenance/yes/2023)

# Terraform provider for [Common Short Domain product](https://github.com/idealo/transport_csd)

The Common Short Domain product gives you cool short domains (AWS Hosted Zones) in your AWS account so you can manage them yourself, without the hassle of a third party.

Currently, we support the following domains where you can get subdomains:

- `idealo.tools`: internal idealo tooling for everyone
- `idealo.com`: idealo components mostly for b2b
- `idealo.de`: idealo components mostly for b2c
- `idealo.co.uk`: idealo components mostly for british b2c
- `idealo.es`: idealo components mostly for spanish b2c
- `idealo.fr`: idealo components mostly for french b2c
- `idealo.it`: idealo components mostly for italian b2c
- `idealo.nl`: idealo components mostly for dutch b2c
- `idealo.pl`: idealo components mostly for polish b2c
- `idealo.pt`: idealo components mostly for portuguese b2c

More domains will follow in future updates. If you're missing one that you need, contact Team Transport.

_Keep in mind that your FQDN shouldn't exceed 64 characters (including the final dot) to retrieve a TLS certificate._

# Installation

## Install from Terraform Registry

You can find our Terraform provider in the [Terraform registry](https://registry.terraform.io/providers/idealo/csd/latest).

Online documentation can also be found [here](https://registry.terraform.io/providers/idealo/csd/latest/docs).

# Usage

```terraform
terraform {
  required_version = "~> 1.3"
  required_providers {
    csd = {
      source  = "idealo/csd"
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
# https://confluence.idealo.cloud/pages/viewpage.action?spaceKey=PTN&title=How+to+authenticate+from+GitHub+to+AWS
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

# Create a Route53 Hosted Zone.
# The lifecycle option prevents Terraform from accidentally removing critical resources.
resource "aws_route53_zone" "shopverwaltung" {
  name = "shopverwaltung.idealo.tools"
  lifecycle {
    prevent_destroy = true
  }
}

# Create zone forwarding in idealo.tools zone via CSD provider
resource "csd_zone" "shopverwaltung" {
  name         = aws_route53_zone.shopverwaltung.name
  name_servers = aws_route53_zone.shopverwaltung.name_servers
}
```

**⚠️ Important:** Keep in mind that the TTL of the NS records for your Hosted Zone can be up to 2 days. So destroying them could lead to extended downtimes for your workloads. We suggest to protect them as shown in the example above and/or separate their automation completely from your product workloads.

# Development

For development notes see [DEV.md](DEV.md).

---

Made with ❤️ and ✨ by [🌐 Team Transport](https://github.com/orgs/idealo/teams/transport).

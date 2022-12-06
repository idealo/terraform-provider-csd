# Terraform provider for [Common Domain Name product](https://github.com/idealo/transport_common_domain)

Terraform provider for the common domain product. With this provider you can register DNS zones under the `idealo.tools` domain. For example `jira.idealo.tools` or `confluence.idealo.tools`.

_Keep in mind that your FQDN shouldn't exceed 64 characters (including the final dot) to retrieve a TLS certificate._

# Installation

_tbd_

# Usage

```terraform
terraform {
  required_version = ">= 1.3"
  required_providers {
    idealo-tools = {
      source  = "github.com/idealo/terraform-provider-idealo-tools"
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

# Setup idealo-tools provider (will use the AWS provider internally)
provider "idealo_tools" {}

# Setup OIDC provider
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
  name = "shopverwaltung.idealo.cloud"
}

# Create zone forwarding in idealo-tools zone
resource "idealo_tools_zone" "shopverwaltung" {
  name         = aws_route53_zone.shopverwaltung.name
  name_servers = aws_route53_zone.shopverwaltung.name_servers
}
```

---

Made with 💖 by Team Transport.

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

# ⚠️ Disclaimer

> With great power comes great responsibility.

Owning your own zone under an idealo.TLD comes with some responsibilities.

## Cookies

Customers log into idealo.de and other idealo TLDs with a cookie that is valid for that domain and its subdomains which includes your hosted zone. This could lead to some unwanted site effects you must be aware of. For example, if you create a CNAME pointing to an external FQDN, the cookie will be readable by that third party. So this external service provider could read that cookie and in the worst case impersonate our customer. From a security perspective, this might be unwanted behaviour. So if you point DNS records to third parties, take care that cookies are not forwarded to them. If you're unsure please contact us or the Security team to clarify how to deal with your specific scenario.

As an example, let's say you serve the wishlist component from you AWS account. For that, you registered the subdomain wishlist.idealo.de with our CSD product. That means that we delegate the zone wishlist.idealo.de to your account. In your account, you then create DNS resource records pointing to the wishlist component, for example an ALB inside your account.
Imagine you use a third party service like Salesforce that requires you to point DNS entries under your hosted zone to their service. For example, a CNAME salesforce.wishlist.idealo.de pointing to service.salesforce.com. This would mean that Salesforce is now able to read the customer's cookie and therefore is able to impersonate that customer. In that case, contact security to make sure that you comply with our security requirements.

## Mail servers

By controlling your own zone, you're also able to set records for your own mail servers. These mail servers would be able to send mails with a sender under subdomain for example wishlist.idealo.de. These mails should be well crafted and aligned with company standards from the design, legal and security departments.

If you plan to set up email communication under your subdomain, you must talk to the mentioned departments first to make you follow the idealo guidelines.

If you have any other questions about your hosted zone setup, feel free to reach out to Team Transport.

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
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "~>4.3"

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
```

## Use case 1: Hosted zone delegation

```terraform
# Create a Route53 Hosted Zone.
# The lifecycle option prevents Terraform from accidentally removing critical resources.
resource "aws_route53_zone" "shopverwaltung" {
  name = "shopverwaltung.idealo.tools"
  lifecycle {
    prevent_destroy = true
  }
}

# Create zone delegation in idealo.tools zone via CSD provider
resource "csd_zone_delegation" "shopverwaltung" {
  name         = aws_route53_zone.shopverwaltung.name
  name_servers = aws_route53_zone.shopverwaltung.name_servers
}
```

**⚠️ Important:** Keep in mind that the TTL of the NS records for your Hosted Zone can be up to 2 days. So destroying them could lead to extended downtimes for your workloads. We suggest to protect them as shown in the example above and/or separate their automation completely from your product workloads.

## Use case 2: Route traffic through Akamai

```terraform
resource "csd_record" "wishlist_idealo_de_cname" {
  name  = "wishlist.idealo.de"
  rrtype  = "cname"
  value = "wishlist.edgekey.net"
  ttl   = 3600
}

resource "csd_record" "_acme_challenge_wishlist_idealo_de_txt" {
  name  = "_acme_challenge.wishlist.idealo.de"
  rrtype  = "txt"
  value = "LeisahxaiQu8ayah2aiwe9Que5saiy4o"
  ttl   = 60
}
```

Follow the detailed documentation on how to setup the Akamai property [here](https://backstage.idealo.tools/catalog/default/component/CSD/docs/#use-case-forward-traffic-to-akamai). If you have any questions about the property, please ask the [SECURITY](https://teams.microsoft.com/l/channel/19%3a77eca9f9ee784e04988b4b8c29814e0b%40thread.tacv2/%25F0%259F%259B%25A1%25EF%25B8%258F%2520PT%2520Security?groupId=424df2ed-7bad-42b5-9c93-2a74f5acd0e1&tenantId=21956b19-fed2-44b7-90cf-b6d281c0a42a) team. They will gladly help you with that.


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

Now follow the proper upgrade procedure described [here](#Upgrade from v1.x to v2.x).

# Development

For development notes see [DEV.md](DEV.md).




---

Made with ❤️ and ✨ by [🌐 Team Transport](https://github.com/orgs/idealo/teams/transport).

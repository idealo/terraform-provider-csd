# Terraform provider for [Common Domain Name product](https://github.com/idealo/transport_common_domain)

Terraform provider for the common domain product. With this provider you can register DNS zones under the `idealo.tools` domain. For example `jira.idealo.tools` or `confluence.idealo.tools`.

_Keep in mind that your FQDN shouldn't exceed 64 characters (including the final dot) to retrieve a TLS certificate for it. Any alternate domain names can be up to 256 characters long._

# Installation

_tbd_

# Usage

```terraform
provider "aws" {}
provider "idealo_tools" {}

resource "aws_route53_zone" "shopverwaltung" {
  name = "shopverwaltung.idealo.cloud"
}

resource "idealo_tools_zone" "shopverwaltung" {
  name         = aws_route53_zone.shopverwaltung.name
  name_servers = aws_route53_zone.shopverwaltung.name_servers
}
```

---

Made with 💖 by Team Transport.

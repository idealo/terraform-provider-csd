terraform {
  required_providers {
    idealo-tools = {
      version = "~> 0.0.0"
      source = "idealo.com/transport/idealo-tools"
    }
    aws = {}
  }
}

provider "idealo-tools" {}
provider "aws" {}

resource "aws_route53_zone" "myzone" {
  name = "myapp.idealo.tools"
}

resource "idealo-tools-zone" "myzone" {
  name = aws_route53_zone.myzone.name
  name_servers = aws_route53_zone.myzone.name_servers
}

output "myapp" {
  value = idealo-tools-zone.myzone.name
}

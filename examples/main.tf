terraform {
  required_providers {
    csd = {
      version = "~> 0.0.1"
      source = "idealo.com/transport/csd"
    }
    #aws = {}
  }
  required_version = "~> 1.0"
}

provider "csd" {}

#provider "aws" {}

#resource "aws_route53_zone" "myzone" {
#  name = "myapp.idealo.tools"
#}
#
#resource "csd_zone" "myzone" {
#  name = aws_route53_zone.myzone.name
#  name_servers = aws_route53_zone.myzone.name_servers
#}

data "csd_zones" "all" {}

output "zones" {
  value = data.csd_zones.all
}

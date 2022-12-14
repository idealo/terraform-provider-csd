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

data "csd_zone" "jira" {
  name = "jira.idealo.tools"
}

output "jira" {
  value = data.csd_zone.jira
}

resource "csd_zone" "confluence" {
  name = "confluence.idealo.tools"
  name_servers = [
    "ns1.aws.example.net",
    "ns2.aws.example.com"
  ]
  owner = "123456789"
}

output "confluence" {
  value = csd_zone.confluence
}

terraform {
  required_providers {
    csd = {
      version = "~> 0.0.1"
      source = "idealo.com/transport/csd"
    }
    aws = {

    }
  }
  required_version = "~> 1.0"
}

provider "csd" {}

provider "aws" {}

#resource "aws_route53_zone" "myzone" {
#  name = "myapp.idealo.tools"
#}
#
#resource "csd_zone" "myzone" {
#  name = aws_route53_zone.myzone.name
#  name_servers = aws_route53_zone.myzone.name_servers
#}

#output "test" {
#  value = csd_zone.myzone.name
#}

data "csd_zones" "all" {}

output "test_data_read_zones" {
  value = data.csd_zones.all
}

#data "csd_zone" "jira" {
#  name = "jira.idealo.tools"
#}
#
#output "test_data_read_zone" {
#  value = data.csd_zone.jira
#}
#
#resource "csd_zone" "confluence" {
#  name = "confluence.idealo.tools"
#  name_servers = [
#    "ns1.aws.example.net",
#    "ns2.aws.example.com"
#  ]
#  owner = "123456789"
#}
#
#output "test_resource_create_zone" {
#  value = csd_zone.confluence
#}

terraform {
  required_version = "~>1.3"
  required_providers {
    csd = {
      version = "~>2.0"
      source  = "idealo/csd"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~>4.8"
    }
  }
}

provider "aws" {
  region              = "eu-central-1"
  allowed_account_ids = ["433744410943"]
}

provider "csd" {}

resource "aws_route53_zone" "my_zone_delegation" {
  name = "myzone.idealo.tools"
}

#resource "csd_zone_delegation" "my_zone_delegation" {
#  name = aws_route53_zone.my_zone_delegation.name
#  name_servers = aws_route53_zone.my_zone_delegation.name_servers
#}
#
#output "test_name" {
#  value = csd_zone_delegation.my_zone_delegation.name
#}
#
#output "test_name_servers" {
#  value = csd_zone_delegation.my_zone_delegation.name_servers
#}

#data "csd_zone_delegations" "all" {}
#
#output "test_data_read_zone_delegations" {
#  value = data.csd_zone_delegations.all
#}
#
#data "csd_zone_delegation" "my_zone_delegation" {
#  name = "myzone.idealo.tools"
#}
#
#output "test_data_read_zone_delegation" {
#  value = data.csd_zone_delegation.my_zone_delegation.name_servers
#}

#resource "csd_zone_delegation" "confluence" {
#  name = "confluence.idealo.tools"
#  name_servers = [
#    "ns1.aws.example.net",
#    "ns2.aws.example.com"
#  ]
#  owner = "123456789"
#}
#
#output "test_resource_create_zone_delegation" {
#  value = csd_zone_delegation.confluence
#}

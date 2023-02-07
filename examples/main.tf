terraform {
  required_version = "~> 1.3"
  required_providers {
    csd = {
      version = "~> 1.0"
      source  = "idealo.com/transport/csd"
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

resource "aws_route53_zone" "myzone" {
  name = "myzone.idealo.tools"
}

#resource "csd_zone" "myzone" {
#  name = aws_route53_zone.myzone.name
#  name_servers = aws_route53_zone.myzone.name_servers
#}
#
#output "test_name" {
#  value = csd_zone.myzone.name
#}
#
#output "test_name_servers" {
#  value = csd_zone.myzone.name_servers
#}

#data "csd_zones" "all" {}
#
#output "test_data_read_zones" {
#  value = data.csd_zones.all
#}
#
#data "csd_zone" "myzone" {
#  name = "myzone.idealo.tools"
#}
#
#output "test_data_read_zone" {
#  value = data.csd_zone.myzone.name_servers
#}

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

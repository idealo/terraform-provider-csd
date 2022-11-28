terraform {
  required_providers {
    idealo-tools = {
      version = "0.0.1"
      source = "idealoo.com/transport/idealo-tools"
    }
    #aws = {}
  }
  required_version = "~> 1.0"
}

provider "idealo-tools" {}
#provider "aws" {}

#resource "aws_route53_zone" "myzone" {
#  name = "myapp.idealo.tools"
#}
#
#resource "idealo_tools_zone" "myzone" {
#  name = aws_route53_zone.myzone.name
#  name_servers = aws_route53_zone.myzone.name_servers
#}


data "idealo_tools_zones" "all" {}


output "jira_zone" {
  value = data.idealo_tools_zones.all.zones
}

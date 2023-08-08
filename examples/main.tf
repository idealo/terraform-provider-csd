terraform {
  required_version = "~>1.3"
  required_providers {
    csd = {
      version = "~>2.0"
      source  = "idealo/csd"
      #source  = "idealo.com/transport/csd"
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
  name = "myzone2.idealo.tools"
}

resource "csd_zone_delegation" "myzone" {
  name         = aws_route53_zone.myzone.name
  name_servers = aws_route53_zone.myzone.name_servers
}

resource "aws_route53_record" "myrecord" {
  zone_id = aws_route53_zone.myzone.id
  name    = "foo"
  type    = "TXT"
  ttl     = 60
  records = ["bar"]
}

resource "csd_record" "myrecord" {
  name   = "_acme-challenge.myrecord.myzone2.idealo.tools"
  rrtype = "TXT"
  value  = "foobar"
}

resource "csd_record" "myrecord2" {
  name   = "myrecord2.myzone2.idealo.tools"
  rrtype = "CNAME"
  value  = "foobar.edgekey.net."
}

#data "csd_record" "myrecord" {
#  name = "myrecord.idealo.tools"
#}

#output "myrecord" {
#  value = data.csd_record.myrecord.value
#}

#data "csd_records" "all" {}

#output "all_records" {
#  value = data.csd_records.all
#}

#resource "aws_route53_zone" "my_zone" {
#  name = "myzone.idealo.tools"
#}

#resource "csd_zone_delegation" "my_zone_delegation" {
#  name = aws_route53_zone.my_zone.name
#  name_servers = aws_route53_zone.my_zone.name_servers
#}

#output "test_name" {
#  value = csd_zone_delegation.my_zone_delegation.name
#}

#output "test_name_servers" {
#  value = csd_zone_delegation.my_zone_delegation.name_servers
#}

#data "csd_zone_delegations" "all" {}

#output "test_data_read_zone_delegations" {
#  value = data.csd_zone_delegations.all
#}

#data "csd_zone_delegation" "my_zone_delegation" {
#  name = "myzone.idealo.tools"
#}

#output "test_data_read_zone_delegation" {
#  value = data.csd_zone_delegation.my_zone_delegation.name_servers
#}


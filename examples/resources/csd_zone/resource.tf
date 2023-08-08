resource "aws_route53_zone" "my_zone" {
  name = "myzone.idealo.tools"
}

resource "csd_zone_delegation" "my_zone_delegation" {
  name         = aws_route53_zone.my_zone.name
  name_servers = aws_route53_zone.my_zone.name_servers
}

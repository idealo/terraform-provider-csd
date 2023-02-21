resource "aws_route53_zone" "myzone" {
  name = "myzone.idealo.tools"
}

resource "csd_zone" "myzone" {
  name         = aws_route53_zone.myzone.name
  name_servers = aws_route53_zone.myzone.name_servers
}

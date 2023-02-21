data "csd_zones" "all" {}

output "all_zones" {
  value = data.csd_zones.all
}

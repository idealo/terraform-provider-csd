data "csd_zone_delegations" "all" {}

output "all_zone_delegations" {
  value = data.csd_zone_delegations.all
}

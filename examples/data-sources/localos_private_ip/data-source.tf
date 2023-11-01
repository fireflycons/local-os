data "localos_private_ip" "my_ips" {}

output "primary_ip" {
  description = "IP of NIC connected to default gateway"
  value       = my_ips.primary.ip
}

output "primary_cidr" {
  description = "/32 CIDR of NIC connected to default gateway"
  value       = my_ips.primary.cidr
}

output "primary_name" {
  description = "Interface name of NIC connected to default gateway"
  value       = my_ips.primary.name
}

output "primary_network" {
  description = "CIDR range of network to which NIC connected to default gateway is part of"
  value       = my_ips.primary.name
}

# "secondaries" is a list of all other local NICs found, not including loopback adapter
# Attributes for each list entry are the same as for "primary"

output "first_secondary_ip" {
  value = my_ips.secondaries[0].ip
}
data "localos_public_ip" "my_ip" {}

output "cidr" {
  value = nonsensitive(data.localos_public_ip.cidr)
}

output "ip" {
  value = nonsensitive(data.localos_public_ip.ip)
}
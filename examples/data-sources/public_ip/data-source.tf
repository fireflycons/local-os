data "localos_public_ip" "my_ip" {}

output "cidr" {
    value = data.localos_public_ip.cidr
    sensitive = true
}
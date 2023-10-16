data "localos_folders" "folders" {}

output "home_directory" {
    value = data.localos_folders.folders.home
}

output "ssh_keys_directory" {
    value = data.localos_folders.folders.ssh
}


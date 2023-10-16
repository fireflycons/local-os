data "localos_info" "os_info" {}

output "os_name" {
    value = data.localos_info.os_info.name
}

output "os_arch" {
    value = data.localos_info.os_info.arch
}

output "os_is_windows" {
    value = data.localos_info.os_info.is_windows
}

output "os_path_var" {
    value = data.localos_info.os_info.environment.PATH
    sensitive = true
}


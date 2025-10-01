data "ztw_network_services_groups" "example" {
  name = "Corporate Custom SSH TCP_10022"
}

output "ztw_network_services_groups" {
  value = data.ztw_network_services_groups.example
}

data "ztc_network_services_groups" "example" {
  name = "Corporate Custom SSH TCP_10022"
}

output "ztc_network_services_groups" {
  value = data.ztc_network_services_groups.example
}

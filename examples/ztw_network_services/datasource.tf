data "ztw_network_service" "example" {
  name = "ICMP_ANY"
}

output "ztw_network_service" {
  value = data.ztw_network_service.example
}

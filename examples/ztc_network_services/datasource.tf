data "ztc_network_service" "example" {
  name = "ICMP_ANY"
}

output "ztc_network_service" {
  value = data.ztc_network_service.example
}

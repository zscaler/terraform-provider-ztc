data "ztc_dns_gateway" "example" {
  name = "Example DNS Gateway"
}

output "ztc_dns_gateway_name" {
  value = data.ztc_dns_gateway.example.name
}

data "ztc_forwarding_gateway" "example" {
  name = "example"
}

output "ztc_forwarding_gateway" {
  value = data.ztc_forwarding_gateway.example
}

data "ztw_forwarding_gateway" "example" {
  name = "example"
}

output "ztw_forwarding_gateway" {
  value = data.ztw_forwarding_gateway.example
}

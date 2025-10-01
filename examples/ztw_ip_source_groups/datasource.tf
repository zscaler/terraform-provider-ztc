data "zia_ip_source_groups" "example" {
  name = "example"
}

output "zia_ip_source_groups" {
  value = data.zia_ip_source_groups.example
}

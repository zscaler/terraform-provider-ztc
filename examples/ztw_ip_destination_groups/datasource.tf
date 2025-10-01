data "ztw_ip_destination_groups" "example" {
  name = "example"
}

output "ztw_ip_destination_groups_example" {
  value = data.ztw_ip_destination_groups.example
}

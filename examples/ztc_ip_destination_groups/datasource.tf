data "ztc_ip_destination_groups" "example" {
  name = "example"
}

output "ztc_ip_destination_groups_example" {
  value = data.ztc_ip_destination_groups.example
}

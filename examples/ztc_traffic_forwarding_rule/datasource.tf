data "ztc_ip_pool_groups" "example" {
  name = "example"
}

output "ztc_ip_pool_groups" {
  value = data.ztc_ip_pool_groups.example
}

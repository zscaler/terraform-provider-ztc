data "ztw_ip_pool_groups" "example" {
  name = "example"
}

output "ztw_ip_pool_groups" {
  value = data.ztw_ip_pool_groups.example
}

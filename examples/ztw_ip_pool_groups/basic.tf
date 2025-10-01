resource "ztw_ip_pool_groups" "example" {
  name        = "Example IP Pool"
  description = "Updated Example IP Pool for testing"
  ip_addresses = [
    "192.168.100.0/24"
  ]
}

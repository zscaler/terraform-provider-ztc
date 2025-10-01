data "ztw_network_service" "example" {
  name = "ICMP_ANY"
}


resource "ztw_network_services_groups" "example" {
  name        = "example"
  description = "example"
  services {
    id = [
      data.ztw_network_service.example.id,
    ]
  }
}

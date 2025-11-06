data "ztc_network_service" "example" {
  name = "ICMP_ANY"
}


resource "ztc_network_services_groups" "example" {
  name        = "example"
  description = "example"
  services {
    id = [
      data.ztc_network_service.example.id,
    ]
  }
}

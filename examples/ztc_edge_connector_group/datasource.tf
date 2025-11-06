data "ztc_edge_connector_group" "this" {
  name = "zs-cc-vpc-096108eb5d9e68d71-ca-central-1a"
}

output "ztc_edge_connector_group" {
  value = data.ztc_edge_connector_group.aws1_vpc_ca
}

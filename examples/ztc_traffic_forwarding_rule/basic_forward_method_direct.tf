resource "ztc_forwarding_gateway" "ztc_gw01" {
  name           = "ZTC_GW01"
  description    = "Example Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ZIA"
}

data "ztc_location_management" "this" {
  name = "AWS-CAN-ca-central-1-vpc-05c7f364cf47c2b93"
}

data "ztc_ip_destination_groups" "this" {
  name = "example"
}

data "zia_ip_source_groups" "this" {
  name = "example"
}

data "ztc_network_service" "this" {
  name = "ICMP_ANY"
}

data "ztc_network_services_groups" "this" {
  name = "Corporate Custom SSH TCP_10022"
}

# NOTE: To retrieve the src_workload_groups ID information, you must leverage the ZIA Terraform Provider at the moment,
# which returns the exact same information ID since the resource is cross-shared between ZIA and ZTC.

data "zia_workload_groups" "this" {
  name = "WORKLOAD_GROUP01"
}

resource "ztc_traffic_forwarding_rule" "this1" {
  name           = "DIRECT_Forwarding_Rule01"
  description    = "DIRECT Forwarding Rule 01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  type           = "EC_RDR"
  forward_method = "DIRECT"
  src_ips        = ["192.168.200.200"]
  dest_addresses = ["192.168.255.1"]
  wan_selection  = "BALANCED_RULE"
  dest_countries = ["CA", "US"]
  src_workload_groups {
    id = [data.zia_workload_groups.this.id]
  }
  nw_service_groups {
    id = [data.ztc_network_services_groups.this.id]
  }
  nw_services {
    id = [data.ztc_network_service.this.id]
  }
  src_ip_groups {
    id = [data.zia_ip_source_groups.this.id]
  }
  dest_ip_groups {
    id = [data.ztc_ip_destination_groups.this.id]
  }
  locations {
    id = [data.ztc_location_management.this.id]
  }
  proxy_gateway {
    id   = ztc_forwarding_gateway.ztc_gw01.id
    name = ztc_forwarding_gateway.ztc_gw01.name
  }
}


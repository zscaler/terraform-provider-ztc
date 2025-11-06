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

data "zia_ip_source_groups" "this" {
  name = "example"
}

# NOTE: To retrieve the src_workload_groups ID information, you must leverage the ZIA Terraform Provider at the moment,
# which returns the exact same information ID since the resource is cross-shared between ZIA and ZTC.

data "zia_workload_groups" "this" {
  name = "WORKLOAD_GROUP01"
}

resource "ztc_traffic_forwarding_rule" "this1" {
  name           = "ECZPA_Forwarding_Rule01"
  description    = "ECZPA Forwarding Rule 01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  type           = "EC_RDR"
  forward_method = "ECZPA"
  src_ips        = ["192.168.200.200"]
  wan_selection  = "BALANCED_RULE"
  src_workload_groups {
    id = [data.zia_workload_groups.this.id]
  }
  src_ip_groups {
    id = [data.zia_ip_source_groups.this.id]
  }
  locations {
    id = [data.ztc_location_management.this.id]
  }
  zpa_application_segments {
    id = [18612387, 18616051] # There is no current way to retrieve the ID information which is required for the configuration.
  }
  zpa_application_segment_groups {
    id = [18612386, 18661615] # There is no current way to retrieve the ID information which is required for the configuration.
  }
}


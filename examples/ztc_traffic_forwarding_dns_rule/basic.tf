## Example Usage - DNS Forwarding Group - Action REDIR_REQ

resource "ztc_dns_forwarding_gateway" "ztc_dns01" {
  name                             = "DNS_FW_GW01"
  ec_dns_gateway_options_primary   = "WAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "WAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_ALLOW_IGNORE_DNAT"
}

resource "ztc_ip_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses   = ["test1.acme.com", "test2.acme.com", "test3.acme.com"]
}

resource "zia_ip_source_groups" "this" {
  name         = "example1"
  description  = "example1"
  ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}

data "ztc_location_management" "this" {
  name = "SJC_01"
}

resource "ztc_traffic_forwarding_dns_rule" "this" {
  name           = "DNS_Rule01"
  description    = "DNS_Rule01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  action         = "REDIR_REQ"
  src_ips        = ["192.168.200.200"]
  dest_addresses = ["server1.acme.com"]
  src_ip_groups {
    id = [zia_ip_source_groups.this.id]
  }
  dest_ip_groups {
    id = [ztc_ip_destination_groups.dstn_fqdn.id]
  }
  locations {
    id = [data.ztc_location_management.this.id]
  }
  dns_gateway {
    id   = ztc_dns_forwarding_gateway.this.id
    name = ztc_dns_forwarding_gateway.this.name
  }
}

## Example Usage - DNS Forwarding Group - Action REDIR_ZPA

resource "ztc_dns_forwarding_gateway" "ztc_dns01" {
  name                             = "DNS_FW_GW01"
  ec_dns_gateway_options_primary   = "WAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "WAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_ALLOW_IGNORE_DNAT"
}

resource "ztc_ip_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses   = ["test1.acme.com", "test2.acme.com", "test3.acme.com"]
}

resource "zia_ip_source_groups" "this" {
  name         = "example1"
  description  = "example1"
  ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}

resource "ztc_ip_pool_groups" "this" {
  name        = "Example IP Pool"
  description = "Updated Example IP Pool for testing"
  ip_addresses = [
    "10.0.0.0/24"
  ]
}

data "ztc_location_management" "this" {
  name = "SJC_01"
}

resource "ztc_traffic_forwarding_dns_rule" "this" {
  name           = "DNS_Rule01"
  description    = "DNS_Rule01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  action         = "REDIR_ZPA"
  src_ips        = ["192.168.200.200"]
  dest_addresses = ["server1.acme.com"]
  src_ip_groups {
    id = [zia_ip_source_groups.this.id]
  }
  locations {
    id = [data.ztc_location_management.this.id]
  }
  zpa_ip_group {
    id   = ztc_ip_pool_groups.this.id
    name = ztc_ip_pool_groups.this.name
  }
}

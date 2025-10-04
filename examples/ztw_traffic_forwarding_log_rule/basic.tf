## Example Usage - DNS Forwarding Group - Action REDIR_REQ

resource "ztw_zia_forwarding_gateway" "this" {
  name             = "ZTW_GW01"
  description      = "Example Forwarding Gateway 1"
  fail_closed      = true
  primary_type     = "MANUAL_OVERRIDE"
  secondary_type   = "MANUAL_OVERRIDE"
  manual_primary   = "1.1.1.1"
  manual_secondary = "2.2.2.2"
  type             = "ZIA"
}

data "ztw_location_management" "this" {
  name = "SJC_01"
}

resource "ztw_traffic_forwarding_log_rule" "this" {
  name           = "Log_Rule01"
  description    = "Log_Rule01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  forward_method = "ECSELF"

  locations {
    id = [data.ztw_location_management.this.id]
  }
  proxy_gateway {
    id   = ztw_zia_forwarding_gateway.this.id
    name = tw_zia_forwarding_gateway.this.name
  }
}

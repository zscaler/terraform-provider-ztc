## Example Usage - ZIA Primary and Secondary Type AUTO

resource "ztc_zia_forwarding_gateway" "ztc_gw01" {
  name           = "ZTC_GW01"
  description    = "Example Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ZIA"
}

## Example Usage - ZIA Primary and Secondary Type DC

resource "ztc_zia_forwarding_gateway" "ztc_gw02" {
  name             = "ZTC_GW02"
  description      = "Example Forwarding Gateway 2"
  fail_closed      = true
  primary_type     = "DC"
  secondary_type   = "DC"
  manual_primary   = "zrh1.svpn.zscalerbeta.net"
  manual_secondary = "syseng.svpn.zscalerbeta.net"
  type             = "ZIA"
}

## Example Usage - ZIA Primary and Secondary Type MANUAL_OVERRIDE

resource "ztc_zia_forwarding_gateway" "ztc_gw03" {
  name             = "ZTC_GW03"
  description      = "Example Forwarding Gateway 3"
  fail_closed      = true
  primary_type     = "MANUAL_OVERRIDE"
  secondary_type   = "MANUAL_OVERRIDE"
  manual_primary   = "1.1.1.1"
  manual_secondary = "2.2.2.2"
  type             = "ZIA"
}

## Example Usage - Log Forwarding Gateway Primary and Secondary Type AUTO

resource "ztc_zia_forwarding_gateway" "ztc_log01" {
  name           = "LOG_FW_GW01"
  description    = "Example Log Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ECSELF"
}

## Example Usage - Log Forwarding Gateway Primary and Secondary Type DC

resource "ztc_zia_forwarding_gateway" "ztc_log02" {
  name             = "LOG_FW_GW02"
  description      = "Example Log Forwarding Gateway 2"
  fail_closed      = true
  primary_type     = "DC"
  secondary_type   = "DC"
  manual_primary   = "zrh1.svpn.zscalerbeta.net"
  manual_secondary = "syseng.svpn.zscalerbeta.net"
  type             = "ECSELF"
}

## Example Usage - Log Forwarding Gateway Primary and Secondary Type MANUAL_OVERRIDE

resource "ztc_zia_forwarding_gateway" "ztc_log03" {
  name             = "LOG_FW_GW03"
  description      = "Example Log Forwarding Gateway 3"
  fail_closed      = true
  primary_type     = "MANUAL_OVERRIDE"
  secondary_type   = "MANUAL_OVERRIDE"
  manual_primary   = "1.1.1.1"
  manual_secondary = "2.2.2.2"
  type             = "ECSELF"
}

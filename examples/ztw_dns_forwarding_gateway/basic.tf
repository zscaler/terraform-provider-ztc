
## Example Usage - DNS Gateway Primary and Secondary DNS Address

resource "ztw_dns_forwarding_gateway" "ztw_dns01" {
  name             = "DNS_FW_GW01"
  primary_ip       = "4.4.4.4"
  secondary_ip     = "8.8.8.8"
  failure_behavior = "FAIL_ALLOW_IGNORE_DNAT"
}


## Example Usage - DNS Gateway ECDNS Options Primary and Secondary

resource "ztw_dns_forwarding_gateway" "ztw_dns02" {
  name                             = "DNS_FW_GW02"
  ec_dns_gateway_options_primary   = "WAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "WAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_ALLOW_IGNORE_DNAT"
}

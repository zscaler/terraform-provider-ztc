resource "ztc_dns_gateway" "example" {
  name                             = "Example DNS Gateway"
  ec_dns_gateway_options_primary   = "LAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "LAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_RET_ERR"
}

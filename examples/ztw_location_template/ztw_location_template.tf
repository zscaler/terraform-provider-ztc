resource "ztw_location_template" "this" {
  name = "testAcc_location_template"
  desc = "Location Template Test"
  template {
    template_prefix           = "testAcc-tf"
    aup_timeout_in_days       = 0
    auth_required             = false
    ips_control               = true
    ofw_enabled               = true
    xff_forward_enabled       = true
    enforce_bandwidth_control = true
    up_bandwidth              = 10
    dn_bandwidth              = 10
  }
}

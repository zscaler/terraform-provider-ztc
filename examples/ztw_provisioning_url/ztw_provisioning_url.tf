resource "ztw_location_template" "this" {
  name = "testAcc_location_template${random_id.test.hex}"
  desc = "Location Template Test"
  template {
    template_prefix           = "testAcc-tf"
    aup_timeout_in_days       = 1
    auth_required             = false
    caution_enabled           = false
    aup_enabled               = true
    ips_control               = true
    ofw_enabled               = true
    xff_forward_enabled       = true
    enforce_bandwidth_control = true
    up_bandwidth              = 10000
    dn_bandwidth              = 10000
  }
}

resource "ztw_provisioning_url" "example" {
  name          = "testAcc_provisioning_url"
  desc          = "Updated Provisioning URL Test"
  prov_url_type = "CLOUD"
  prov_url_data {
    form_factor         = "SMALL"
    cloud_provider_type = "AWS"
    location_template {
      id = ztw_location_template.this.id
    }
  }
}

output "name" {
  value = ztw_provisioning_url.example.prov_url
}

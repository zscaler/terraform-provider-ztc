data "ztw_supported_regions" "this" {
  name = "US_EAST_1"
}


resource "ztw_public_cloud_info" "this" {
  name       = "AWSAccount01"
  cloud_type = "AWS"
  account_details {
    aws_account_id           = "123456789"
    aws_role_name            = "zscaler-role"
    cloud_watch_group_arn    = "DISABLED"
    event_bus_name           = "zscaler-bus-123456-zscalerthree.net"
    trouble_shooting_logging = true
    trusted_account_id       = "123456789"
    trusted_role             = "arn:aws:iam::123456789:role/ZscalerTagDiscoveryRole"
  }
  supported_regions {
    id = [data.ztw_supported_regions.this.id]
  }
}

data "ztw_location_template" "this" {
  name = "aws-ca-central-1"
}

output "ztw_location_template" {
  value = data.ztw_location_template.this
}

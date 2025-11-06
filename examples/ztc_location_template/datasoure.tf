data "ztc_location_template" "this" {
  name = "aws-ca-central-1"
}

output "ztc_location_template" {
  value = data.ztc_location_template.this
}

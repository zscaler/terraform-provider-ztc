data "ztc_activation_status" "this" {}

output "ztc_activation_status" {
  value = data.ztc_activation_status.this
}

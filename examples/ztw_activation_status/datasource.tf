data "ztw_activation_status" "this" {}

output "ztw_activation_status" {
  value = data.ztw_activation_status.this
}

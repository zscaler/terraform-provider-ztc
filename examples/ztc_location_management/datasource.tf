data "ztc_location_management" "aws_vpc_05c7f364cf47c2b93" {
  name = "AWS-CAN-ca-central-1-vpc-05c7f364cf47c2b93"
}

output "ztc_location_management" {
  value = data.ztc_location_management.aws_vpc_05c7f364cf47c2b93
}

data "aws_availability_zones" "current" {}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.18.1"

  name = "ptf-backend-vpc"
  cidr = var.vpc_cidr

  azs = slice(data.aws_availability_zones.current.names, 0, 2)
  public_subnets = [
    cidrsubnet(var.vpc_cidr, 4, 0),
    cidrsubnet(var.vpc_cidr, 4, 1)
  ]
  private_subnets = [
    cidrsubnet(var.vpc_cidr, 4, 2),
    cidrsubnet(var.vpc_cidr, 4, 3)
  ]

  enable_nat_gateway     = var.enable_nat_gateway
  single_nat_gateway     = var.enable_nat_gateway ? (var.one_nat_gateway_per_az ? false : true) : false
  one_nat_gateway_per_az = var.enable_nat_gateway ? var.one_nat_gateway_per_az : false

  enable_dns_hostnames = true # default value
  enable_dns_support   = true # default value
}
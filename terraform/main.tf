locals {
  app_name = "ptf-golang-app"
}

module "vpc" {
  app_name               = local.app_name
  source                 = "./usecases/vpc"
  vpc_cidr               = "10.0.0.0/16"
  enable_nat_gateway     = true
  one_nat_gateway_per_az = false
}

module "ecr" {
  source   = "./modules/ecr"
  app_name = local.app_name
}

module "ssm_parameters" {
  source   = "./modules/ssm"
  app_name = local.app_name
  secrets = [
    "TEST_SECRET",
  ]
}

module "ecs" {
  source   = "./usecases/ecs"
  app_name = local.app_name
}
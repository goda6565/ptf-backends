module "vpc" {
  source                 = "./usecases/vpc"
  vpc_cidr               = "10.0.0.0/16"
  enable_nat_gateway     = true
  one_nat_gateway_per_az = false
}

module "ecr" {
  source   = "./modules/ecr"
  app_name = "ptf-golang-app"
}
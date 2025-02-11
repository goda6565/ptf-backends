variable "vpc_cidr" {
  type        = string
  description = "The CIDR block for the VPC"
}

variable "enable_nat_gateway" {
  type        = bool
  description = "Enable NAT Gateway for private subnets"
}

variable "one_nat_gateway_per_az" {
  type        = bool
  description = "Enable one NAT Gateway per Availability Zone"
  default     = false
}
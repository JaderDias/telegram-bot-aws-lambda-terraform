
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = var.name
  cidr = var.vpc_cidr

  azs            = [for az in ["a", "b", "c"] : "${var.aws_region}${az}"]
  public_subnets = var.public_subnet_cidrs

  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_nat_gateway = false # if true, costs $0.73/day
  single_nat_gateway = false

  create_igw = true
}
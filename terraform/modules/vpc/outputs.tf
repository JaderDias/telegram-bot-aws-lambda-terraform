output "vpc_id" {
  value = module.vpc.vpc_id
}

output "public_subnets" {
  value = module.vpc.public_subnets
}

output "sg_for_lambda" {
  value = data.aws_security_group.default
}

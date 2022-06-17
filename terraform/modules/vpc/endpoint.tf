module "vpc_endpoints" {
  source = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"

  vpc_id             = module.vpc.vpc_id
  security_group_ids = [data.aws_security_group.default.id]

  endpoints = {
    ssm = {
      service             = "ssm"
      private_dns_enabled = true
      subnet_ids          = module.vpc.public_subnets
      security_group_ids  = [aws_security_group.vpc_tls.id]
    },
    kms = {
      service             = "kms"
      private_dns_enabled = true
      subnet_ids          = module.vpc.public_subnets
      security_group_ids  = [aws_security_group.vpc_tls.id]
    },
  }

  tags = {
    Project  = "Secret"
    Endpoint = "true"
  }
}
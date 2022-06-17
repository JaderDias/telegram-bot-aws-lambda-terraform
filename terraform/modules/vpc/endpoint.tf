module "vpc_endpoints" {
  source = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"

  vpc_id             = aws_default_vpc.default.id
  security_group_ids = [aws_default_security_group.default.id]

  endpoints = {
    ssm = {
      service             = "ssm"
      private_dns_enabled = true
      subnet_ids          = [for subnet in aws_default_subnet.default_subnet : subnet.id]
      security_group_ids  = [aws_security_group.vpc_tls.id]
    },
    kms = {
      service             = "kms"
      private_dns_enabled = true
      subnet_ids          = [for subnet in aws_default_subnet.default_subnet : subnet.id]
      security_group_ids  = [aws_security_group.vpc_tls.id]
    },
  }

  tags = {
    Project  = "Secret"
    Endpoint = "true"
  }
}
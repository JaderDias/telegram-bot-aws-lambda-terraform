output "vpc_id" {
  value = aws_default_vpc.default.id
}

output "public_subnets" {
  value = aws_default_subnet.default_subnet
}

output "sg_for_lambda" {
  value = aws_default_security_group.default.id
}

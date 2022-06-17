data "aws_internet_gateway" "default" {
  filter {
    name   = "attachment.vpc-id"
    values = [aws_default_vpc.default.id]
  }
}

resource "aws_default_route_table" "region_route_table" {
  default_route_table_id = aws_default_vpc.default.default_route_table_id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = data.aws_internet_gateway.default.id
  }
  route {
    ipv6_cidr_block = "::/0"
    gateway_id = data.aws_internet_gateway.default.id
  }
}
resource "aws_default_vpc" "default" {
}

resource "aws_default_subnet" "default_subnet" {
  for_each          = toset(["a", "b", "c"])
  availability_zone = "${var.aws_region}${each.key}"
}
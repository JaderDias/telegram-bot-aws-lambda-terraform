variable "aws_region" {
  description = "AWS region for all resources."
  type        = string
}

variable "name" {
  description = "the name fo the vpc"
  default     = "lambda-vpc"
}

variable "vpc_cidr" {
  description = "cidr range for the vpc"
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "cidrs for public subnets"
  default     = ["10.0.96.0/20", "10.0.112.0/20", "10.0.128.0/20"]
}
variable "aws_region" {
  description = "AWS region for all resources."
  type        = string
}

variable "name" {
  description = "the name fo the vpc"
  default     = "lambda-vpc"
}

variable "tags" {
  description = "values for tags"
  type        = map(string)
}
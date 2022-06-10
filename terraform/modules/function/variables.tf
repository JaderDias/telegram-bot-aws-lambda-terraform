variable "aws_ssm_parameter_arn" {
  type    = string
  default = null
}
variable "function_name" {
  type = string
}
variable "source_dir" {
  type = string
}
variable "lambda_handler" {
  type = string
}
variable "schedule_expression" {
  type    = string
  default = null
}
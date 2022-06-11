variable "aws_ssm_parameter_arn" {
  type    = string
  default = null
}
variable "function_name" {
  type = string
}
variable "lambda_handler" {
  type = string
}
variable "language" {
  type = string
}
variable "source_dir" {
  type = string
}
variable "s3_bucket_arn" {
  type = string
}
variable "s3_bucket_id" {
  type = string
}
variable "schedule_expression" {
  type    = string
  default = null
}
variable "bucket_arn" {
  type    = string
  default = null
}
variable "function_name" {
  type = string
}
variable "source_file" {
  type = string
}
variable "lambda_handler" {
  type = string
}
variable "schedule_expression" {
  type    = string
  default = null
}
variable "secret_arn" {
  type    = string
  default = null
}
variable "create_url" {
  type    = bool
  default = false
}
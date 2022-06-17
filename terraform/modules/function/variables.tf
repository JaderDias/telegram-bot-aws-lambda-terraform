variable "function_name" {
  type = string
}
variable "lambda_handler" {
  type = string
}
variable "language" {
  type    = string
  default = null
}
variable "source_dir" {
  type = string
}
variable "schedule_expression" {
  type    = string
  default = null
}
variable "url_authorization_type" {
  type    = string
  default = "AWS_IAM"
}
variable "ssm_key_arn" {
  type    = string
  default = null
}
variable "ssm_parameter" {
}

variable "subnets" {
  description = "subnets for lambda function"
}

variable "security_groups" {
  description = "security group ids for lambda function"
  type        = list(string)
  default     = []
}

variable "efs_access_point_arn" {
  description = "efs access point arn"
}

variable "local_mount_path" {
  description = "local mount path in lambda function. must start with '/mnt/'"
  default     = "/mnt/efs"
}

variable "efs_mount_targets" {
  description = "efs file system mount targets"
}
variable "tags" {
  description = "tags for lambda function"
  type        = map(string)
}
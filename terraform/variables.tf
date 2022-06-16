variable "aws_region" {
  description = "AWS region for all resources."
  type    = string
}

variable "telegram_bot_tokens" {
  description = "Telegram bot tokens."

  type      = string
  sensitive = true
}

variable "name" {
  description = "application name"
  default     = "telegram-bot"
}

variable "efs_throughput_mode" {
  description = "Throughput mode for the file system. Defaults to bursting. Valid values: bursting, provisioned. When using provisioned, also set provisioned_throughput_in_mibps"
  default     = null
}

variable "efs_provisioned_throughput" {
  description = "The throughput, measured in MiB/s, that you want to provision for the file system. Only applicable with `throughput_mode` set to `provisioned`"
  default     = null
}
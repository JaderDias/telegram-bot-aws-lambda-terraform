variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "eu-central-1"
}

variable "telegram_bot_token" {
  description = "Telegram bot token."

  type      = string
  sensitive = true
}
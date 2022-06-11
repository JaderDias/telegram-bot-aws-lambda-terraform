variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "eu-central-1"
}

variable "telegram_bot_tokens" {
  description = "Telegram bot tokens."

  type      = string
  sensitive = true
}
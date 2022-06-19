variable "aws_region" {
  description = "AWS region for all resources."
  type    = string
}

variable "telegram_bot_tokens" {
  description = "Telegram bot tokens."

  type      = string
  sensitive = true
}
resource "aws_kms_key" "aws_ssm_key" {
  description = "ssm telegram bot token key"
}

resource "aws_ssm_parameter" "telegram_bot_tokens" {
  name   = "${terraform.workspace}_telegram_bot_tokens"
  type   = "SecureString"
  value  = var.telegram_bot_tokens
  key_id = aws_kms_key.aws_ssm_key.arn
  tags = {
    environment = terraform.workspace,
    deployment  = random_pet.this.id,
  }
}
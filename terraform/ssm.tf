resource "aws_ssm_parameter" "telegram_bot_tokens" {
  name  = "${terraform.workspace}_telegram_bot_tokens"
  type  = "SecureString"
  value = var.telegram_bot_tokens
  tags = {
    environment = terraform.workspace,
    deployment  = random_pet.this.id,
  }
}

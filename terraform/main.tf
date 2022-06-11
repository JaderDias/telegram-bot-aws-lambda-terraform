terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.17.1"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

resource "random_pet" "this" {
  length = 2
}

resource "aws_ssm_parameter" "telegram_bot_tokens" {
  name  = "telegram_bot_tokens"
  type  = "SecureString"
  value = var.telegram_bot_tokens
}

resource "aws_s3_bucket" "bucket" {
  bucket = "my-bucket-${random_pet.this.id}"
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.bucket.id
  acl    = "private"
}

#module "send_message_function" {
#  source = "./modules/function"

#  function_name       = "send_message_function-${random_pet.this.id}"
#  lambda_handler      = "send_message"
#  source_dir          = "../bin/send_message"
#  schedule_expression = "rate(60 minutes)"
#}

module "reply_function" {
  for_each = jsondecode(nonsensitive(var.telegram_bot_tokens))
  source   = "./modules/function"

  aws_ssm_parameter_arn = aws_ssm_parameter.telegram_bot_tokens.arn
  function_name         = "reply_${each.key}_${random_pet.this.id}"
  lambda_handler        = "reply"
  language              = each.key
  source_dir            = "../bin/reply"
  s3_bucket_arn         = aws_s3_bucket.bucket.arn
  s3_bucket_id          = aws_s3_bucket.bucket.id
}

resource "null_resource" "register_webhook" {
  triggers = {
    always_run = "${timestamp()}"
  }
  provisioner "local-exec" {
    working_dir = "../golang/register"
    command = format("go run . '%s'", jsonencode({
      for k, v in module.reply_function : k => v["function_url"]
    }))
    interpreter = ["bash", "-c"]
  }
}
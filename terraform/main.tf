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

resource "aws_ssm_parameter" "telegram_bot_token" {
  name  = "telegram_bot_token"
  type  = "SecureString"
  value = var.telegram_bot_token
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
  source = "./modules/function"

  aws_ssm_parameter_arn = aws_ssm_parameter.telegram_bot_token.arn
  function_name         = "reply_function-${random_pet.this.id}"
  lambda_handler        = "reply"
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
    command     = "go run . ${var.telegram_bot_token} ${module.reply_function.function_url}"
    interpreter = ["bash", "-c"]
  }
  depends_on = [
    module.reply_function
  ]
}
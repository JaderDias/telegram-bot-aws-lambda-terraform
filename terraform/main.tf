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

resource "null_resource" "upload_language" {
  triggers = {
    aws_s3_bucket = aws_s3_bucket.bucket.id
    nl_update     = filesha1("../nl.csv"),
    sh_update     = filesha1("../sh.csv")
  }
  provisioner "local-exec" {
    working_dir = "../golang/upload"
    command = format(
      "go run . %s '%s'",
      var.aws_region,
      aws_s3_bucket.bucket.id
    )
    interpreter = ["bash", "-c"]
  }
}

module "send_function" {
  source              = "./modules/function"
  function_name       = "${terraform.workspace}_send_${random_pet.this.id}"
  lambda_handler      = "send"
  source_dir          = "../bin/send"
  schedule_expression = "rate(60 minutes)"
  s3_bucket_arn       = aws_s3_bucket.bucket.arn
  s3_bucket_id        = aws_s3_bucket.bucket.id
  ssm_parameter_arn   = aws_ssm_parameter.telegram_bot_tokens.arn
  ssm_parameter_name  = aws_ssm_parameter.telegram_bot_tokens.name
}

module "reply_function" {
  for_each = jsondecode(nonsensitive(var.telegram_bot_tokens))
  source   = "./modules/function"

  function_name      = "${terraform.workspace}_reply_${each.key}_${random_pet.this.id}"
  lambda_handler     = "reply"
  language           = each.key
  source_dir         = "../bin/reply"
  s3_bucket_arn      = aws_s3_bucket.bucket.arn
  s3_bucket_id       = aws_s3_bucket.bucket.id
  ssm_parameter_arn  = aws_ssm_parameter.telegram_bot_tokens.arn
  ssm_parameter_name = aws_ssm_parameter.telegram_bot_tokens.name
}

resource "null_resource" "register_webhook" {
  triggers = {
    always_run = "${timestamp()}"
  }
  provisioner "local-exec" {
    working_dir = "../golang/register"
    command = format(
      "go run . %s %s '%s'",
      var.aws_region,
      aws_ssm_parameter.telegram_bot_tokens.name,
      jsonencode({
        for k, v in module.reply_function : k => v["function_url"]
      }),
    )
    interpreter = ["bash", "-c"]
  }
}

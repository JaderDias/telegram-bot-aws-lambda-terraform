terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
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

#resource "null_resource" "generate_dictionary" {
#  triggers = {
#    always_run = "${timestamp()}"
#  }
#  provisioner "local-exec" {
#    command     = "../generate_dictionary.sh ../golang/reply"
#    interpreter = ["bash", "-c"]
#  }
#}

#module "send_message_function" {
#  source = "./modules/function"

#  function_name       = "send_message_function-${random_pet.this.id}"
#  lambda_handler      = "send_message"
#  source_file         = "../bin/send_message"
#  schedule_expression = "rate(60 minutes)"
#}

module "reply_function" {
  source = "./modules/function"

  function_name  = "reply_function-${random_pet.this.id}"
  lambda_handler = "reply"
  source_file    = "../bin/reply"
  secret_arn     = aws_ssm_parameter.telegram_bot_token.arn
#  depends_on = [
#    resource.null_resource.generate_dictionary
  #]
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
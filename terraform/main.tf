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

module "vpc" {
  source = "./modules/vpc"

  name                 = "${terraform.workspace}-${var.name}-vpc"
  aws_region           = var.aws_region
  vpc_cidr             = "10.0.0.0/16"
  public_subnet_cidrs  = ["10.0.96.0/20", "10.0.112.0/20", "10.0.128.0/20"]
  private_subnet_cidrs = ["10.0.0.0/20", "10.0.16.0/20", "10.0.32.0/20"]
}

module "efs" {
  source = "./modules/efs"

  name                   = "${terraform.workspace}-${var.name}-efs"
  subnet_ids             = module.vpc.private_subnets
  security_group_ids     = [module.vpc.sg_for_lambda.id]
  provisioned_throughput = var.efs_provisioned_throughput
  throughput_mode        = var.efs_throughput_mode
}

module "upload_function" {
  source = "./modules/function"

  function_name     = "${terraform.workspace}_upload_${random_pet.this.id}"
  lambda_handler    = "upload"
  source_dir        = "../bin/upload"
  subnet_ids        = module.vpc.private_subnets
  security_groups   = [module.vpc.sg_for_lambda.id]
  ssm_parameter_arn = aws_ssm_parameter.telegram_bot_tokens.arn
  ssm_key_arn       = aws_kms_key.aws_ssm_key.arn

  efs_access_point_arn = module.efs.access_point_arn
  efs_mount_targets    = module.efs.mount_targets
}

resource "aws_lambda_invocation" "upload" {
  function_name = module.upload_function.function_name
  input         = ""
  triggers = {
    nl_update = filesha1("../nl.csv"),
    sh_update = filesha1("../sh.csv")
  }
  depends_on = [
    module.upload_function,
  ]
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = "../bin/reply"
  output_path = "../bin/reply.zip"
}

module "reply_function" {
  for_each = jsondecode(nonsensitive(var.telegram_bot_tokens))
  source   = "./modules/function"

  function_name      = "${terraform.workspace}_reply_${each.key}_${random_pet.this.id}"
  lambda_handler     = "reply"
  language           = each.key
  source_dir         = "../bin/reply"
  subnet_ids         = module.vpc.private_subnets
  security_groups    = [module.vpc.sg_for_lambda.id]
  ssm_parameter_arn  = aws_ssm_parameter.telegram_bot_tokens.arn
  ssm_parameter_name = aws_ssm_parameter.telegram_bot_tokens.name
  ssm_key_arn        = aws_kms_key.aws_ssm_key.arn

  efs_access_point_arn   = module.efs.access_point_arn
  efs_mount_targets      = module.efs.mount_targets
  url_authorization_type = "NONE"
  depends_on = [
    aws_lambda_invocation.upload,
    data.archive_file.lambda_zip
  ]
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
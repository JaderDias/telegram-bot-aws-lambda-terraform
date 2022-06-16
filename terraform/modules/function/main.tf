data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = var.source_dir
  output_path = "${var.source_dir}.zip"
}

resource "aws_lambda_function" "myfunc" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = var.function_name
  role             = aws_iam_role.iam_for_terraform_lambda.arn
  handler          = var.lambda_handler
  source_code_hash = filebase64sha256(data.archive_file.lambda_zip.output_path)
  runtime          = "go1.x"
  timeout          = 30
  environment {
    variables = {
      language             = var.language
      s3_bucket_id         = var.s3_bucket_id
      token_parameter_name = var.ssm_parameter_name
    }
  }
  tags = {
    environment = terraform.workspace
  }
}

resource "aws_lambda_function_url" "url1" {
  function_name      = aws_lambda_function.myfunc.function_name
  qualifier          = ""
  authorization_type = "NONE"

  cors {
    allow_credentials = true
    allow_origins     = ["*"]
    allow_methods     = ["POST"]
    allow_headers     = ["date", "keep-alive"]
    expose_headers    = ["keep-alive", "date"]
    max_age           = 86400
  }
  depends_on = [
    aws_lambda_function.myfunc
  ]
}
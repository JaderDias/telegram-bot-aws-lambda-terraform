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
      token_parameter_name = var.ssm_parameter_name
      LOCAL_MOUNT_PATH     = var.local_mount_path
    }
  }
  tags = {
    environment = terraform.workspace
  }

  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = var.security_groups
  }


  file_system_config {
    arn              = var.efs_access_point_arn
    local_mount_path = var.local_mount_path
  }

  # Explicitly declare dependency on EFS mount target.
  # When creating or updating Lambda functions, mount target must be in 'available' lifecycle state.
  depends_on = [
    var.efs_mount_targets,
    data.archive_file.lambda_zip
  ]
}
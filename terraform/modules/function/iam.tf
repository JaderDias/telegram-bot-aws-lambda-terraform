data "aws_kms_alias" "aws_ssm_key" {
  name = "alias/aws/ssm"
}

data "aws_iam_policy_document" "lambda_exec_role_policy" {
  statement {
    actions = [
      "ssm:GetParameter",
    ]
    resources = [
      var.ssm_parameter_arn
    ]
  }
  statement {
    actions = [
      "kms:Decrypt",
    ]
    resources = [
      data.aws_kms_alias.aws_ssm_key.arn
    ]
  }
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    effect    = "Allow"
    resources = ["arn:aws:logs:*:*:*"]
  }
  statement {
    actions = [
      "s3:GetObject",
      "s3:PutObject",
    ]
    effect = "Allow"
    resources = [
      "${var.s3_bucket_arn}/*",
      var.s3_bucket_arn,
    ]
  }
}

# Lambda function policy
resource "aws_iam_policy" "lambda_policy" {
  name        = "${var.function_name}-lambda-policy"
  description = "${var.function_name}-lambda-policy"
  policy      = data.aws_iam_policy_document.lambda_exec_role_policy.json
  tags = {
    environment = terraform.workspace
  }
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    effect = "Allow"
  }
}

# Lambda function role
resource "aws_iam_role" "iam_for_terraform_lambda" {
  name               = "${var.function_name}-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
  tags = {
    environment = terraform.workspace
  }
}

# Role to Policy attachment
resource "aws_iam_role_policy_attachment" "terraform_lambda_iam_policy_basic_execution" {
  role       = aws_iam_role.iam_for_terraform_lambda.id
  policy_arn = aws_iam_policy.lambda_policy.arn
}
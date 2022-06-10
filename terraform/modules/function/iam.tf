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
      var.ssm_key_arn
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

data "aws_partition" "current" {}

data "aws_iam_policy" "AmazonElasticFileSystemClientFullAccess" {
  arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonElasticFileSystemClientFullAccess"
}

data "aws_iam_policy" "AWSLambdaVPCAccessExecutionRole" {
  arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy_attachment" "AmazonElasticFileSystemClientFullAccess-attach" {
  role       = aws_iam_role.iam_for_terraform_lambda.id
  policy_arn = data.aws_iam_policy.AmazonElasticFileSystemClientFullAccess.arn
}

resource "aws_iam_role_policy_attachment" "AWSLambdaVPCAccessExecutionRole-attach" {
  role       = aws_iam_role.iam_for_terraform_lambda.id
  policy_arn = data.aws_iam_policy.AWSLambdaVPCAccessExecutionRole.arn
}
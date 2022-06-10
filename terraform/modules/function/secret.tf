data "aws_iam_policy_document" "lambda_exec_role_policy" {
  statement {
    actions = [
      "ssm:GetParameter",
    ]
    resources = [
      var.secret_arn
    ]
  }
}
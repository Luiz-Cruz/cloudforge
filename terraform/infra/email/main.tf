resource "aws_ses_email_identity" "email" {
  email = var.email
}

data "aws_iam_policy_document" "policy" {
  statement {
    actions   = ["SES:SendEmail", "SES:SendRawEmail"]
    resources = [aws_ses_email_identity.email.arn]

    principals {
      identifiers = ["*"]
      type        = "AWS"
    }
  }
}

resource "aws_ses_identity_policy" "policy" {
  identity = aws_ses_email_identity.email.arn
  name     = "cloudforge-ses-policy"
  policy   = data.aws_iam_policy_document.policy.json
}

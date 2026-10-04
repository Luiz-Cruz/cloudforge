variable "environment" {}
variable "lambda_exec_role_arn" {}
variable "sqs_main_queue_arn" {}

# Dummy zip just for initial TF apply without the real Go binary
resource "archive_file" "dummy_lambda" {
  type        = "zip"
  output_path = "${path.module}/dummy.zip"
  source {
    content  = "dummy"
    filename = "bootstrap"
  }
}

# API Lambda Function
resource "aws_lambda_function" "api" {
  function_name = "cloudforge-api-${var.environment}"
  role          = var.lambda_exec_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2"
  filename      = archive_file.dummy_lambda.output_path

  environment {
    variables = {
      CLOUD       = "AWS"
      APPLICATION = "API"
    }
  }
}

# Worker Lambda Function
resource "aws_lambda_function" "worker" {
  function_name = "cloudforge-worker-${var.environment}"
  role          = var.lambda_exec_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2"
  filename      = archive_file.dummy_lambda.output_path

  environment {
    variables = {
      CLOUD       = "AWS"
      APPLICATION = "JOB"
    }
  }
}

# SQS Trigger for Worker
resource "aws_lambda_event_source_mapping" "sqs_trigger" {
  event_source_arn = var.sqs_main_queue_arn
  function_name    = aws_lambda_function.worker.arn
  batch_size       = 10
}

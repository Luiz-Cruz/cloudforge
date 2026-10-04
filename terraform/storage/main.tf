variable "environment" {}

# DynamoDB Table to store the SAGA State
resource "aws_dynamodb_table" "saga_state" {
  name           = "cloudforge-saga-state-${var.environment}"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "transaction_id"

  attribute {
    name = "transaction_id"
    type = "S"
  }
}

# SQS Dead Letter Queue
resource "aws_sqs_queue" "dlq" {
  name = "cloudforge-dlq-${var.environment}"
}

# Main SQS Queue with Redrive Policy
resource "aws_sqs_queue" "main_queue" {
  name = "cloudforge-queue-${var.environment}"
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = 3
  })
}

output "main_queue_arn" {
  value = aws_sqs_queue.main_queue.arn
}

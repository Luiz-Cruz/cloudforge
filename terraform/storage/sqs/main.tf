resource "aws_sqs_queue" "dlq" {
  name                      = var.dlq_name
  message_retention_seconds = 1209600 # 14 days
}

resource "aws_sqs_queue" "primary" {
  name                       = var.queue_name
  visibility_timeout_seconds = 900    # 15 minutes (aligned with SAGA execution timeout)
  message_retention_seconds  = 345600 # 4 days

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = var.max_receive_count
  })
}

resource "aws_sqs_queue_redrive_allow_policy" "dlq_allow" {
  queue_url = aws_sqs_queue.dlq.id

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue"
    sourceQueueArns   = [aws_sqs_queue.primary.arn]
  })
}

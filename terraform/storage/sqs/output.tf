output "queue_arn" {
  description = "ARN of the primary provisioning SQS queue"
  value       = aws_sqs_queue.primary.arn
}

output "queue_url" {
  description = "URL of the primary provisioning SQS queue"
  value       = aws_sqs_queue.primary.url
}

output "dlq_arn" {
  description = "ARN of the Dead Letter Queue"
  value       = aws_sqs_queue.dlq.arn
}

output "dlq_url" {
  description = "URL of the Dead Letter Queue"
  value       = aws_sqs_queue.dlq.url
}

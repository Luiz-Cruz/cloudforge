output "lambda_arn" {
  description = "ARN of the SAGA Worker Lambda function"
  value       = aws_lambda_function.lambda_job.arn
}

output "topic_arn" {
  description = "ARN of the notifications SNS topic"
  value       = aws_sns_topic.notifications_topic.arn
}

output "alarms_topic_arn" {
  description = "ARN of the alarms SNS topic"
  value       = aws_sns_topic.alarms_topic.arn
}

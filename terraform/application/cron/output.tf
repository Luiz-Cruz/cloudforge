output "lambda_arn" {
  description = "ARN of the Cron Maintenance Lambda function"
  value       = aws_lambda_function.lambda_cron.arn
}

output "event_rule_arn" {
  description = "ARN of the EventBridge rule"
  value       = aws_cloudwatch_event_rule.maintenance_schedule.arn
}

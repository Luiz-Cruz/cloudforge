output "api_gateway_url" {
  description = "Production API Gateway URL"
  value       = module.api.api_url
}

output "dynamodb_table_name" {
  description = "Production DynamoDB SAGA table name"
  value       = module.dynamodb.table_name
}

output "s3_bucket_name" {
  description = "Production S3 bucket name"
  value       = module.s3.bucket_name
}

output "sqs_queue_url" {
  description = "Production SQS queue URL"
  value       = module.sqs.queue_url
}

output "sqs_dlq_url" {
  description = "Production Dead Letter Queue URL"
  value       = module.sqs.dlq_url
}

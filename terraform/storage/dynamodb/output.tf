output "table_arn" {
  description = "ARN of the DynamoDB SAGA table"
  value       = aws_dynamodb_table.saga_state.arn
}

output "table_name" {
  description = "Name of the DynamoDB SAGA table"
  value       = aws_dynamodb_table.saga_state.name
}

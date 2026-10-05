output "api_url" {
  description = "Invoke URL for the API Gateway stage"
  value       = aws_api_gateway_stage.api_stage.invoke_url
}

output "lambda_arn" {
  description = "ARN of the API Lambda function"
  value       = aws_lambda_function.lambda_api.arn
}

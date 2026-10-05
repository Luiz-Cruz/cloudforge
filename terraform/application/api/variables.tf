variable "lambda_api_name" {
  description = "Name of the API Lambda function"
  type        = string
  default     = "cloudforge-api"
}

variable "lambda_role_arn" {
  description = "IAM role ARN for the Lambda function"
  type        = string
}

variable "image_uri" {
  description = "Container image URI for Lambda execution"
  type        = string
}

variable "env_var" {
  description = "Environment variables for the Lambda function"
  type        = map(string)
  default     = {}
}

variable "gateway_name" {
  description = "Name of the API Gateway"
  type        = string
  default     = "cloudforge-gateway"
}

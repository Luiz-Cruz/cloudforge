variable "lambda_cron_name" {
  description = "Name of the Cron Maintenance Lambda function"
  type        = string
  default     = "cloudforge-cron"
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

variable "schedule_expression" {
  description = "EventBridge cron expression for maintenance execution"
  type        = string
  default     = "cron(0 2 * * ? *)" # Everyday at 02:00 AM UTC
}

variable "lambda_job_name" {
  description = "Name of the SAGA Worker Lambda function"
  type        = string
  default     = "cloudforge-job"
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

variable "sqs_arn" {
  description = "ARN of the primary SQS queue to trigger the worker"
  type        = string
}

variable "topic_name" {
  description = "Name of the SNS topic for notifications"
  type        = string
  default     = "cloudforge-notifications"
}

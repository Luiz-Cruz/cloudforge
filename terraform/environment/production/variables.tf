variable "aws_region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-east-1"
}

variable "ecr_repository_name" {
  description = "Name of the ECR repository holding the unified container image"
  type        = string
  default     = "cloudforge-container"
}

variable "ecr_image_tag" {
  description = "Image tag to deploy"
  type        = string
  default     = "latest"
}

variable "table_name" {
  description = "Name of the DynamoDB table"
  type        = string
  default     = "cloudforge-saga-state-prod"
}

variable "bucket_name" {
  description = "Name of the S3 storage bucket"
  type        = string
  default     = "cloudforge-artifacts-production"
}

variable "role_name" {
  description = "Base name for IAM role"
  type        = string
  default     = "cloudforge-lambda-prod"
}

variable "lambda_api_name" {
  description = "Name of the API Lambda function"
  type        = string
  default     = "cloudforge-api-prod"
}

variable "lambda_job_name" {
  description = "Name of the SAGA Worker Lambda function"
  type        = string
  default     = "cloudforge-job-prod"
}

variable "lambda_cron_name" {
  description = "Name of the Cron Maintenance Lambda function"
  type        = string
  default     = "cloudforge-cron-prod"
}

variable "gateway_name" {
  description = "Name of the API Gateway"
  type        = string
  default     = "cloudforge-api-gateway-prod"
}

variable "email_sender" {
  description = "Verified sender email in Amazon SES"
  type        = string
  default     = "no-reply@cloudforge.io"
}

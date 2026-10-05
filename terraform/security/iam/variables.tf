variable "role_name" {
  description = "Base name for IAM role and policy"
  type        = string
  default     = "cloudforge-lambda"
}

variable "dynamodb_tables" {
  description = "List of DynamoDB table ARNs"
  type        = list(string)
}

variable "bucket_arn" {
  description = "ARN of the S3 storage bucket"
  type        = string
}

variable "sqs_queues" {
  description = "List of SQS queue ARNs"
  type        = list(string)
  default     = []
}

variable "sns_topics" {
  description = "List of SNS topic ARNs"
  type        = list(string)
  default     = ["*"]
}

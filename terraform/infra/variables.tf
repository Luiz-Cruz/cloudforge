variable "aws_region" {
  description = "AWS region for infrastructure resources"
  type        = string
  default     = "us-east-1"
}

variable "repository_name" {
  description = "ECR image repository name"
  type        = string
  default     = "cloudforge-container"
}

variable "email_sender" {
  description = "SES verified email address"
  type        = string
  default     = "no-reply@cloudforge.io"
}

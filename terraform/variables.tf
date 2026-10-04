variable "aws_region" {
  default = "us-east-1"
}

variable "environment" {
  description = "The environment name (e.g., local, prod)"
  type        = string
}

variable "is_local" {
  description = "Flag to determine if we are targeting LocalStack"
  type        = bool
  default     = false
}

variable "table_name" {
  description = "Name of the DynamoDB table for CloudForge SAGA state"
  type        = string
  default     = "cloudforge-saga-state"
}

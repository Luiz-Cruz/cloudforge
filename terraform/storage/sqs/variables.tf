variable "queue_name" {
  description = "Name of the primary SQS provisioning queue"
  type        = string
  default     = "cloudforge-provisioning-queue"
}

variable "dlq_name" {
  description = "Name of the Dead Letter Queue (DLQ)"
  type        = string
  default     = "cloudforge-provisioning-dlq"
}

variable "max_receive_count" {
  description = "Number of times a message is delivered to the source queue before being moved to the dead-letter queue"
  type        = number
  default     = 3
}

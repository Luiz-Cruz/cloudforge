data "aws_caller_identity" "current" {}

data "aws_ecr_repository" "repository" {
  name = var.ecr_repository_name
}

data "aws_ecr_image" "lambda_image" {
  repository_name = var.ecr_repository_name
  image_tag       = var.ecr_image_tag
}

locals {
  image_uri = "${data.aws_ecr_repository.repository.repository_url}@${data.aws_ecr_image.lambda_image.id}"
}

module "dynamodb" {
  source     = "../../storage/dynamodb"
  table_name = var.table_name
}

module "s3" {
  source      = "../../storage/s3"
  bucket_name = var.bucket_name
}

module "sqs" {
  source     = "../../storage/sqs"
  queue_name = "cloudforge-provisioning-queue-prod"
  dlq_name   = "cloudforge-provisioning-dlq-prod"
}

module "iam" {
  source          = "../../security/iam"
  role_name       = var.role_name
  bucket_arn      = module.s3.bucket_arn
  dynamodb_tables = [module.dynamodb.table_arn]
  sqs_queues      = [module.sqs.queue_arn, module.sqs.dlq_arn]
  sns_topics      = [module.job.topic_arn, module.job.alarms_topic_arn]
}

module "api" {
  source          = "../../application/api"
  lambda_api_name = var.lambda_api_name
  lambda_role_arn = module.iam.lambda_role_arn
  image_uri       = local.image_uri
  gateway_name    = var.gateway_name

  env_var = {
    "SERVER"                 = "AWS"
    "CLOUD"                  = "AWS"
    "APPLICATION"            = "API"
    "TABLE_NAME"             = module.dynamodb.table_name
    "QUEUE_URL"              = module.sqs.queue_url
    "STORAGE_BUCKET"         = module.s3.bucket_name
    "EMAIL_SENDER"           = var.email_sender
    "NOTIFICATION_TOPIC_ARN" = module.job.topic_arn
  }
}

module "job" {
  source          = "../../application/job"
  lambda_job_name = var.lambda_job_name
  lambda_role_arn = module.iam.lambda_role_arn
  image_uri       = local.image_uri
  sqs_arn         = module.sqs.queue_arn
  topic_name      = "cloudforge-notifications-prod"

  env_var = {
    "SERVER"                 = "AWS"
    "CLOUD"                  = "AWS"
    "APPLICATION"            = "JOB"
    "TABLE_NAME"             = module.dynamodb.table_name
    "QUEUE_URL"              = module.sqs.queue_url
    "DLQ_URL"                = module.sqs.dlq_url
    "STORAGE_BUCKET"         = module.s3.bucket_name
    "EMAIL_SENDER"           = var.email_sender
    "NOTIFICATION_TOPIC_ARN" = module.job.topic_arn
  }
}

module "cron" {
  source           = "../../application/cron"
  lambda_cron_name = var.lambda_cron_name
  lambda_role_arn  = module.iam.lambda_role_arn
  image_uri        = local.image_uri

  env_var = {
    "SERVER"                 = "AWS"
    "CLOUD"                  = "AWS"
    "APPLICATION"            = "CRON"
    "TABLE_NAME"             = module.dynamodb.table_name
    "STORAGE_BUCKET"         = module.s3.bucket_name
    "EMAIL_SENDER"           = var.email_sender
    "NOTIFICATION_TOPIC_ARN" = module.job.topic_arn
  }
}

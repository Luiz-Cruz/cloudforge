terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region                      = var.aws_region
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = var.is_local
  skip_metadata_api_check     = var.is_local
  skip_requesting_account_id  = var.is_local

  dynamic "endpoints" {
    for_each = var.is_local ? [1] : []
    content {
      dynamodb       = "http://localhost:4566"
      s3             = "http://localhost:4566"
      sns            = "http://localhost:4566"
      sqs            = "http://localhost:4566"
      eventbridge    = "http://localhost:4566"
      apigateway     = "http://localhost:4566"
      lambda         = "http://localhost:4566"
      stepfunctions  = "http://localhost:4566"
    }
  }
}

module "storage" {
  source = "./storage"
  environment = var.environment
}

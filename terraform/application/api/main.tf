data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_lambda_function" "lambda_api" {
  function_name = var.lambda_api_name
  role          = var.lambda_role_arn
  timeout       = 30
  memory_size   = 256
  image_uri     = var.image_uri
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = var.env_var
  }
}

locals {
  swagger_template = templatefile("${path.module}/swagger-terraform.json", {
    aws_region     = data.aws_region.current.region
    aws_account_id = data.aws_caller_identity.current.account_id
    lambda_id      = aws_lambda_function.lambda_api.function_name
  })
}

resource "aws_api_gateway_rest_api" "api_gateway" {
  name        = var.gateway_name
  description = "CloudForge API Gateway"
  body        = local.swagger_template

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id

  triggers = {
    redeployment = sha1(aws_api_gateway_rest_api.api_gateway.body)
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "api_stage" {
  depends_on    = [aws_api_gateway_deployment.api_deployment]
  deployment_id = aws_api_gateway_deployment.api_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  stage_name    = "prod"
}

resource "aws_lambda_permission" "apigw_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_api.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api_gateway.execution_arn}/*/*/*"
}

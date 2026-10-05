resource "aws_lambda_function" "lambda_cron" {
  function_name = var.lambda_cron_name
  role          = var.lambda_role_arn
  timeout       = 300
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

resource "aws_cloudwatch_event_rule" "maintenance_schedule" {
  name                = "${var.lambda_cron_name}-schedule"
  description         = "Triggers CloudForge maintenance cleanup for stale failed environments"
  schedule_expression = var.schedule_expression
}

resource "aws_cloudwatch_event_target" "maintenance_target" {
  rule      = aws_cloudwatch_event_rule.maintenance_schedule.name
  target_id = "MaintenanceLambdaTarget"
  arn       = aws_lambda_function.lambda_cron.arn

  input = jsonencode({
    name = "CLEANUP_FAILED_ENVIRONMENTS"
  })
}

resource "aws_lambda_permission" "eventbridge_permission" {
  statement_id  = "AllowEventBridgeInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_cron.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.maintenance_schedule.arn
}

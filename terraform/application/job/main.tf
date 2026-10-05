resource "aws_sns_topic" "notifications_topic" {
  name = var.topic_name
}

resource "aws_sns_topic" "alarms_topic" {
  name = "${var.lambda_job_name}-alarms"
}

resource "aws_lambda_function" "lambda_job" {
  function_name = var.lambda_job_name
  role          = var.lambda_role_arn
  timeout       = 300
  memory_size   = 512
  image_uri     = var.image_uri
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = merge(var.env_var, {
      "NOTIFICATION_TOPIC_ARN" = aws_sns_topic.notifications_topic.arn
    })
  }
}

resource "aws_lambda_event_source_mapping" "sqs_trigger" {
  event_source_arn = var.sqs_arn
  function_name    = aws_lambda_function.lambda_job.arn
  batch_size       = 5
  enabled          = true
}

resource "aws_cloudwatch_metric_alarm" "job_errors" {
  alarm_name          = "${var.lambda_job_name}-errors"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = 300
  statistic           = "Sum"
  threshold           = 1
  alarm_description   = "Triggers when CloudForge SAGA provisioning encounters an unhandled error"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alarms_topic.arn]

  dimensions = {
    FunctionName = aws_lambda_function.lambda_job.function_name
  }
}

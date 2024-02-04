resource "aws_apigatewayv2_api" "api" {
  name          = "Terraform HTTP API Example"
  protocol_type = "HTTP"
}

resource "aws_cloudwatch_log_group" "logs" {
  name = "/aws/vendedlogs/tf_http_logs"
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id      = aws_apigatewayv2_api.api.id
  auto_deploy = true
  name        = "$default"
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.logs.arn
    format = jsonencode({"requestId":"$context.requestId", "ip":"$context.identity.sourceIp", "requestTime":"$context.requestTime", "httpMethod":"$context.httpMethod","routeKey":"$context.routeKey", "status":"$context.status","protocol":"$context.protocol", "responseLength":"$context.responseLength", "integrationError":"$context.integrationErrorMessage" })
  }
}

# #######################################
# ## Open endpoint                     ##
# #######################################

resource "aws_apigatewayv2_integration" "open_integration" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  integration_uri        = module.lambda_function_responder.lambda_function_invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_open" {
  api_id             = aws_apigatewayv2_api.api.id
  target             = "integrations/${aws_apigatewayv2_integration.open_integration.id}"
  route_key          = "GET /open"
  operation_name     = "get_open_operation"
  authorization_type = "NONE"
}

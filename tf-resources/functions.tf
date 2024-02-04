module "lambda_function_responder" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 6.0"

  timeout             = 300
  source_path         = "../lessons/"
  function_name       = "http_lessons"
  handler             = "main"
  runtime             = "go1.x"
  create_sam_metadata = true
  publish             = true

  allowed_triggers = {
    APIGatewayAny = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
    }
  }
}

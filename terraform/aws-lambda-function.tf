resource "archive_file" "lambda" {
  source_dir  = local.lambdas_path
  output_path = "files/${local.lambda_local_name}.zip"
  type        = "zip"
}

resource "aws_lambda_function" "lambda" {

  filename      = archive_file.lambda.output_path
  function_name = local.lambda_local_name
  role          = aws_iam_role.lambda-role.arn
  handler       = "app"

  description = "Lambda com a função de criar um indice"

  source_code_hash = archive_file.lambda.output_base64sha256

  runtime = var.lambda_runtime

#  environment {
#    variables = {
#      REGION = var.aws_region
#      SURVEY_URL = var.survey_url
#      CONVERSATION_URL = var.conversation_url
#      SURVEY_RESPONSE_COLLECTOR_TABLE = local.survey-response-collector-table
#      CALLBACK_REQUEST_TABLE = local.callback-request-table
#      CONVERSATION_TABLE = local.conversation-table
#    }
#  }


  vpc_config {
    # Every subnet should be able to reach an EFS mount target in the same Availability Zone. Cross-AZ mounts are not permitted.
    subnet_ids         = ["subnet-68ab3159","subnet-2e0b5d48","subnet-1c606812","subnet-893a69d6","subnet-6bea0c27","subnet-a6b9e487"]
    security_group_ids = ["sg-00329077bffb598ff"]
  }

  tags = local.common_tags
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "1.15.0"
    }

    docker = {
      source  = "kreuzwerker/docker"
      version = "2.17.0"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-west-2"
}

locals {
  postgres_identifier    = "gbd"
  postgres_name          = "gbd"
  postgres_user_name     = "gbd"
  postgres_user_password = "SlipperyBeef"
  postgres_instance_name = "gbd"
  postgres_db_password   = "SlipperyBeef"
  postgres_port          = 5432
}

resource "aws_cognito_user_pool" "bco_pool" {
  name = "dev-bco-pool"

  admin_create_user_config {
    allow_admin_create_user_only = true
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }
}

resource "aws_cognito_user_pool_client" "bco_pool_client" {
  name                                 = "dev-bco-pool-client"
  user_pool_id                         = aws_cognito_user_pool.bco_pool.id
  generate_secret                      = false
  allowed_oauth_flows_user_pool_client = false
  supported_identity_providers         = ["COGNITO"]

  explicit_auth_flows = [
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH"
  ]
}

resource "aws_cognito_user" "cognito_user_elliot" {
  user_pool_id       = aws_cognito_user_pool.bco_pool.id
  username           = "elliot"
  temporary_password = "Password1234!"

  attributes = {
    email          = "mpbyrne@gmail.com"
    email_verified = "true"
  }
}

resource "aws_cognito_user" "cognito_user_gofreescout" {
  user_pool_id       = aws_cognito_user_pool.bco_pool.id
  username           = "gofreescout"
  temporary_password = "Password1234!"

  attributes = {
    email          = "mpbyrne@gmail.com"
    email_verified = "true"
  }
}

resource "aws_s3_bucket" "bco_images" {
  bucket = "dev-bco-images"
}

resource "aws_s3_bucket_ownership_controls" "bco_images_ownership" {
  bucket = aws_s3_bucket.bco_images.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "bco_images_acl" {
  depends_on = [
    aws_s3_bucket_ownership_controls.bco_images_ownership, aws_s3_bucket_public_access_block.bco_images_access
  ]

  bucket = aws_s3_bucket.bco_images.id
  acl    = "public-read"
}

resource "aws_s3_bucket_public_access_block" "bco_images_access" {
  bucket = aws_s3_bucket.bco_images.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

data "aws_iam_policy_document" "bco_images_policy" {
  policy_id = "dev-bco-images-policy"

  statement {
    actions = [
      "s3:GetObject"
    ]
    effect = "Allow"
    resources = [
      "${aws_s3_bucket.bco_images.arn}/*"
    ]
    principals {
      type        = "*"
      identifiers = ["*"]
    }
    sid = "S3BCOPublicAccess"
  }
}

resource "aws_s3_bucket_policy" "bco_images_bucket_policy" {
  bucket = aws_s3_bucket.bco_images.id
  policy = data.aws_iam_policy_document.bco_images_policy.json
}

resource "aws_s3_bucket" "bco_images_private" {
  bucket = "dev-bco-images-private"
}

# data "aws_iam_policy_document" "bco_images_private_policy" {
#   policy_id = "dev-bco-images-private-policy"
#
#   statement {
#     actions = [
#       "s3:GetObject"
#     ]
#     effect    = "Allow"
#     resources = [
#       aws_s3_bucket.bco_images_private.arn
#     ]
#     principals {
#       type        = "*"
#       identifiers = ["*"]
#     }
#   }
# }
#
# resource "aws_s3_bucket_policy" "bco_images_private_bucket_policy" {
#   bucket = aws_s3_bucket.bco_images_private.id
#   policy = data.aws_iam_policy_document.bco_images_private_policy.json
# }

resource "aws_s3_bucket" "bco_functions" {
  bucket = "dev-bco-functions"
}

resource "aws_lambda_function" "bco_images_relay" {
  function_name = "dev-bco-images-relay"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = 512
  timeout       = 15
  role          = aws_iam_role.bco_images_relay.arn
  s3_bucket     = aws_s3_bucket.bco_functions.bucket
  s3_key        = "bootstrap4.zip"

  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.bco_images_private.id
    }
  }
}

resource "aws_iam_role" "bco_images_relay" {
  name = "dev-bco-images-relay-role"
  assume_role_policy = jsonencode(
    {
      Statement = [
        {
          Action = "sts:AssumeRole"
          Effect = "Allow"
          Principal = {
            Service = "lambda.amazonaws.com"
          }
        },
      ]
      Version = "2012-10-17"
    }
  )

  inline_policy {
    name = "dev-bco-images-relay-policy"
    policy = jsonencode(
      {
        Statement = [
          {
            Action   = "s3:GetObject"
            Effect   = "Allow"
            Resource = "${aws_s3_bucket.bco_images_private.arn}/*"
          },
        ]
        Version = "2012-10-17"
      }
    )
  }
}

resource "aws_apigatewayv2_api" "bco_images_relay" {
  name                       = "dev-bco-images-relay"
  protocol_type              = "HTTP"
  route_selection_expression = "$request.method $request.path"
}

resource "aws_apigatewayv2_integration" "bco_images_relay" {
  api_id                 = aws_apigatewayv2_api.bco_images_relay.id
  integration_type       = "AWS_PROXY"
  connection_type        = "INTERNET"
  payload_format_version = "1.0"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.bco_images_relay.invoke_arn
  passthrough_behavior   = "WHEN_NO_MATCH"
}

resource "aws_lambda_permission" "bco_images_relay" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bco_images_relay.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.bco_images_relay.execution_arn}/*"
}

resource "aws_apigatewayv2_route" "bco_images_relay" {
  api_id             = aws_apigatewayv2_api.bco_images_relay.id
  route_key          = "$default"
  target             = "integrations/${aws_apigatewayv2_integration.bco_images_relay.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.bco_images_relay.id
}

resource "aws_apigatewayv2_stage" "bco_images_relay" {
  api_id      = aws_apigatewayv2_api.bco_images_relay.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_authorizer" "bco_images_relay" {
  api_id          = aws_apigatewayv2_api.bco_images_relay.id
  name            = "dev-bco-images-relay-authorizer"
  authorizer_type = "JWT"
  identity_sources = [
    "$request.header.Authorization"
  ]
  jwt_configuration {
    audience = [aws_cognito_user_pool_client.bco_pool_client.id]
    issuer   = "https://${aws_cognito_user_pool.bco_pool.endpoint}"
  }
}

provider "postgresql" {
  host            = aws_db_instance.gbd_postgres.address
  port            = local.postgres_port
  database        = local.postgres_name
  username        = local.postgres_user_name
  password        = local.postgres_user_password
  sslmode         = "require"
  connect_timeout = 15
  superuser       = true
}

resource "aws_security_group" "gbd_security_group" {
  name = "gbd-security-group"

  ingress {
    from_port   = local.postgres_port
    to_port     = local.postgres_port
    protocol    = "tcp"
    description = "PostgreSQL"
    cidr_blocks = ["0.0.0.0/0"] // >
  }

  ingress {
    from_port        = local.postgres_port
    to_port          = local.postgres_port
    protocol         = "tcp"
    description      = "PostgreSQL"
    ipv6_cidr_blocks = ["::/0"] // >
  }
}

resource "aws_db_instance" "gbd_postgres" {
  allocated_storage      = 20
  storage_type           = "gp2"
  engine                 = "postgres"
  engine_version         = "15.6"
  instance_class         = "db.m5.large"
  identifier             = local.postgres_identifier
  username               = local.postgres_user_name
  password               = local.postgres_db_password
  publicly_accessible    = true
  vpc_security_group_ids = [aws_security_group.gbd_security_group.id]
  skip_final_snapshot    = true
}

# deployment

data "aws_ecr_authorization_token" "go_server" {
  registry_id = aws_ecr_repository.go_server.registry_id
}

provider "docker" {
  registry_auth {
    address  = split("/", local.ecr_url)[0]
    username = data.aws_ecr_authorization_token.go_server.user_name
    password = data.aws_ecr_authorization_token.go_server.password
  }
}

module "oidc_github" {
  source  = "unfunco/oidc-github/aws"
  version = "1.7.1"

  github_repositories = [
    "michaelByrne/goBoard"
  ]

  iam_role_inline_policies = {
    "actions" : data.aws_iam_policy_document.actions.json
  }
}

data "aws_iam_policy_document" "actions" {
  statement {
    actions = [
      "s3:GetObject",
      "ec2:TerminateInstances",
      "iam:PassRole",
    ]
    effect    = "Allow"
    resources = ["*"]
  }
}


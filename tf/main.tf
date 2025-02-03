terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-west-2"
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

# resource "aws_cognito_user" "cognito_user_abbsworth" {
#   user_pool_id       = aws_cognito_user_pool.bco_pool.id
#   username           = "abbsworth"
#   temporary_password = "Password1234!"
#
#   attributes = {
#     email          = "abigailruthe@gmail.com"
#     email_verified = "true"
#   }
# }

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

# deployment
module "oidc_github" {
  source              = "unfunco/oidc-github/aws"
  version             = "1.7.1"
  attach_admin_policy = true

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
      "ec2:RunInstances",
    ]
    effect    = "Allow"
    resources = ["*"]
  }
}


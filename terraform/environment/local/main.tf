module "main" {
  source      = "../../"
  aws_region  = "us-east-1"
  environment = "local"
  is_local    = true
}

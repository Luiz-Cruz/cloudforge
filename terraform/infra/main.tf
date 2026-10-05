module "ecr" {
  source          = "./image_repository"
  repository_name = var.repository_name
}

module "email" {
  source = "./email"
  email  = var.email_sender
}

output "ecr_repository_url" {
  value = module.ecr.repository_url
}

output "ses_identity_arn" {
  value = module.email.identity_arn
}

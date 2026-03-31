output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "eks_cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
}

output "eks_cluster_endpoint" {
  description = "EKS cluster API server endpoint"
  value       = module.eks.cluster_endpoint
  sensitive   = true
}

output "kubeconfig_command" {
  description = "Run this command to configure kubectl for the cluster"
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
}

output "ecr_frontend_url" {
  description = "ECR repository URL for the frontend image"
  value       = module.ecr.frontend_repository_url
}

output "ecr_backend_url" {
  description = "ECR repository URL for the backend image"
  value       = module.ecr.backend_repository_url
}

output "rds_endpoint" {
  description = "RDS instance hostname (without port) — only reachable within the VPC"
  value       = module.rds.db_endpoint
  sensitive   = true
}

output "rds_username" {
  description = "Master username for the RDS instance"
  value       = module.rds.db_username
}

output "rds_db_name" {
  description = "Database name"
  value       = module.rds.db_name
}

output "docker_login_command" {
  description = "Authenticate Docker with ECR before building/pushing images"
  value       = "aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"
}

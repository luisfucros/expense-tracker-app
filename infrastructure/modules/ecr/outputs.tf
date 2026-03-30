output "frontend_repository_url" {
  description = "ECR repository URL for the frontend image"
  value       = aws_ecr_repository.frontend.repository_url
}

output "backend_repository_url" {
  description = "ECR repository URL for the backend image"
  value       = aws_ecr_repository.backend.repository_url
}

output "frontend_repository_name" {
  value = aws_ecr_repository.frontend.name
}

output "backend_repository_name" {
  value = aws_ecr_repository.backend.name
}

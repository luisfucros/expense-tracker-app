output "db_endpoint" {
  description = "RDS instance hostname (without port) — reachable only within the VPC"
  value       = aws_db_instance.this.address
}

output "db_port" {
  description = "RDS instance port"
  value       = aws_db_instance.this.port
}

output "security_group_id" {
  description = "Security group ID attached to the RDS instance"
  value       = aws_security_group.rds.id
}

output "db_username" {
  description = "Master username for the RDS instance"
  value       = aws_db_instance.this.username
}

output "db_name" {
  description = "Database name"
  value       = aws_db_instance.this.db_name
}

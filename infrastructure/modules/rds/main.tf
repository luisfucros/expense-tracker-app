locals {
  name = "${var.project_name}-${var.environment}"
}

# ── Security Group — only EKS nodes can reach port 3306 ──────────────────────
resource "aws_security_group" "rds" {
  name        = "${local.name}-rds-sg"
  description = "RDS MySQL inbound 3306 from EKS nodes only, no public access"
  vpc_id      = var.vpc_id

  ingress {
    description     = "MySQL from EKS worker nodes"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [var.eks_node_sg_id]
  }

  tags = { Name = "${local.name}-rds-sg" }
}

# ── DB Subnet Group ───────────────────────────────────────────────────────────
resource "aws_db_subnet_group" "this" {
  name        = "${local.name}-db-subnet-group"
  description = "DB subnet group for ${local.name}"
  subnet_ids  = var.db_subnet_ids

  tags = { Name = "${local.name}-db-subnet-group" }
}

# ── Parameter Group ───────────────────────────────────────────────────────────
resource "aws_db_parameter_group" "this" {
  name        = "${local.name}-mysql8"
  family      = "mysql8.0"
  description = "Custom MySQL 8.0 parameters for ${local.name}"

  parameter {
    name  = "character_set_server"
    value = "utf8mb4"
  }

  parameter {
    name  = "collation_server"
    value = "utf8mb4_unicode_ci"
  }

  # Limit connections; increase if connection pool grows
  parameter {
    name  = "max_connections"
    value = "100"
  }

  tags = { Name = "${local.name}-mysql8-params" }
}

# ── RDS Instance ──────────────────────────────────────────────────────────────
resource "aws_db_instance" "this" {
  identifier = "${local.name}-mysql"

  engine         = "mysql"
  engine_version = "8.0"
  instance_class = var.instance_class

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  allocated_storage     = var.allocated_storage
  max_allocated_storage = 50 # Autoscaling ceiling
  storage_type          = "gp3"
  storage_encrypted     = true

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.this.name

  # Never expose the database to the internet
  publicly_accessible = false

  # Single-AZ is fine for this lightweight app.
  # Set multi_az = true for a production workload that requires HA.
  multi_az = false

  backup_retention_period = 7
  backup_window           = "03:00-04:00"
  maintenance_window      = "Mon:04:00-Mon:05:00"
  auto_minor_version_upgrade = true

  # Set deletion_protection = true and skip_final_snapshot = false for production
  deletion_protection       = false
  skip_final_snapshot       = true

  tags = { Name = "${local.name}-mysql" }
}

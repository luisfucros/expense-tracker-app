variable "project_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "db_subnet_ids" {
  type        = list(string)
  description = "Subnet IDs for the RDS subnet group (isolated, no internet route)"
}

variable "eks_node_sg_id" {
  type        = string
  description = "Security group ID of EKS worker nodes — the only source allowed to reach 3306"
}

variable "db_name" {
  type    = string
  default = "expense_tracker"
}

variable "db_username" {
  type    = string
  default = "app"
}

variable "instance_class" {
  type    = string
  default = "db.t3.micro"
}

variable "allocated_storage" {
  type    = number
  default = 20
}

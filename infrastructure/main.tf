data "aws_caller_identity" "current" {}

# ── VPC ───────────────────────────────────────────────────────────────────────
module "vpc" {
  source = "./modules/vpc"

  project_name       = var.project_name
  environment        = var.environment
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
}

# ── ECR ───────────────────────────────────────────────────────────────────────
module "ecr" {
  source = "./modules/ecr"

  project_name = var.project_name
  environment  = var.environment
}

# ── EKS ───────────────────────────────────────────────────────────────────────
module "eks" {
  source = "./modules/eks"

  project_name        = var.project_name
  environment         = var.environment
  aws_region          = var.aws_region
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  eks_cluster_version = var.eks_cluster_version
  node_instance_type  = var.eks_node_instance_type
  node_min_size       = var.eks_node_min_size
  node_max_size       = var.eks_node_max_size
  node_desired_size   = var.eks_node_desired_size
}

# ── RDS ───────────────────────────────────────────────────────────────────────
module "rds" {
  source = "./modules/rds"

  project_name      = var.project_name
  environment       = var.environment
  vpc_id            = module.vpc.vpc_id
  db_subnet_ids     = module.vpc.db_subnet_ids
  eks_node_sg_id    = module.eks.node_security_group_id
  db_name           = var.rds_db_name
  db_username       = var.rds_username
  instance_class    = var.rds_instance_class
  allocated_storage = var.rds_allocated_storage
}

# ── AWS Load Balancer Controller ──────────────────────────────────────────────
# Watches for Ingress objects and creates ALBs automatically.
resource "helm_release" "aws_load_balancer_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"
  version    = "1.7.1"

  set {
    name  = "clusterName"
    value = module.eks.cluster_name
  }

  set {
    name  = "serviceAccount.create"
    value = "true"
  }

  set {
    name  = "serviceAccount.name"
    value = "aws-load-balancer-controller"
  }

  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = module.eks.lbc_iam_role_arn
  }

  set {
    name  = "region"
    value = var.aws_region
  }

  set {
    name  = "vpcId"
    value = module.vpc.vpc_id
  }

  depends_on = [module.eks]
}

# ── Metrics Server ────────────────────────────────────────────────────────────
# Required for HorizontalPodAutoscaler (CPU/memory-based scaling).
resource "helm_release" "metrics_server" {
  name       = "metrics-server"
  repository = "https://kubernetes-sigs.github.io/metrics-server/"
  chart      = "metrics-server"
  namespace  = "kube-system"
  version    = "3.11.0"

  depends_on = [module.eks]
}

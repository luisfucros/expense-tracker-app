locals {
  name = "${var.project_name}-${var.environment}"

  # Auto-compute subnet CIDRs from the VPC CIDR
  public_subnet_cidrs  = [for i, _ in var.availability_zones : cidrsubnet(var.vpc_cidr, 8, i + 1)]
  private_subnet_cidrs = [for i, _ in var.availability_zones : cidrsubnet(var.vpc_cidr, 8, i + 10)]
  db_subnet_cidrs      = [for i, _ in var.availability_zones : cidrsubnet(var.vpc_cidr, 8, i + 20)]
}

# ── VPC ───────────────────────────────────────────────────────────────────────
resource "aws_vpc" "this" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = { Name = "${local.name}-vpc" }
}

# ── Internet Gateway ──────────────────────────────────────────────────────────
resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id

  tags = { Name = "${local.name}-igw" }
}

# ── Public Subnets (ALB) ──────────────────────────────────────────────────────
resource "aws_subnet" "public" {
  count = length(var.availability_zones)

  vpc_id                  = aws_vpc.this.id
  cidr_block              = local.public_subnet_cidrs[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name                                                            = "${local.name}-public-${var.availability_zones[count.index]}"
    "kubernetes.io/role/elb"                                        = "1"
    "kubernetes.io/cluster/${var.project_name}-${var.environment}"  = "shared"
  }
}

# ── Private Subnets (EKS Nodes) ───────────────────────────────────────────────
resource "aws_subnet" "private" {
  count = length(var.availability_zones)

  vpc_id            = aws_vpc.this.id
  cidr_block        = local.private_subnet_cidrs[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = {
    Name                                                            = "${local.name}-private-${var.availability_zones[count.index]}"
    "kubernetes.io/role/internal-elb"                               = "1"
    "kubernetes.io/cluster/${var.project_name}-${var.environment}"  = "shared"
  }
}

# ── Database Subnets (RDS) ────────────────────────────────────────────────────
resource "aws_subnet" "db" {
  count = length(var.availability_zones)

  vpc_id            = aws_vpc.this.id
  cidr_block        = local.db_subnet_cidrs[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = {
    Name = "${local.name}-db-${var.availability_zones[count.index]}"
  }
}

# ── NAT Gateway (single, for cost efficiency) ─────────────────────────────────
# For higher availability, create one NAT Gateway per AZ by iterating over
# public subnets. For this lightweight app a single NAT is sufficient.
resource "aws_eip" "nat" {
  domain = "vpc"

  tags = { Name = "${local.name}-nat-eip" }

  depends_on = [aws_internet_gateway.this]
}

resource "aws_nat_gateway" "this" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public[0].id

  tags = { Name = "${local.name}-nat" }

  depends_on = [aws_internet_gateway.this]
}

# ── Route Tables ──────────────────────────────────────────────────────────────
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }

  tags = { Name = "${local.name}-public-rt" }
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.this.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.this.id
  }

  tags = { Name = "${local.name}-private-rt" }
}

# DB subnets have no outbound internet route — RDS only needs inbound from nodes
resource "aws_route_table" "db" {
  vpc_id = aws_vpc.this.id

  tags = { Name = "${local.name}-db-rt" }
}

# ── Route Table Associations ──────────────────────────────────────────────────
resource "aws_route_table_association" "public" {
  count          = length(aws_subnet.public)
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count          = length(aws_subnet.private)
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "db" {
  count          = length(aws_subnet.db)
  subnet_id      = aws_subnet.db[count.index].id
  route_table_id = aws_route_table.db.id
}

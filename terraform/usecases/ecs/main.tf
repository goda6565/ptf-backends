# VPCとサブネットの情報を取得
locals {
  vpc_name = "${var.app_name}-vpc"
}

data "aws_vpc" "this" {
  filter {
    name   = "tag:Name"
    values = [local.vpc_name]
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "tag:Name"
    values = ["${var.app_name}-public-*"]
  }
}

data "aws_subnets" "private" {
  filter {
    name   = "tag:Name"
    values = ["${var.app_name}-private-*"]
  }
}


# セキュリティーグループの作成
resource "aws_security_group" "alb" {
  name   = "${var.app_name}-alb"
  vpc_id = data.aws_vpc.this.id
}

resource "aws_security_group" "ecs_instance" {
  name   = "${var.app_name}-ecs-instance"
  vpc_id = data.aws_vpc.this.id
}

resource "aws_vpc_security_group_ingress_rule" "alb_from_http" {
  ip_protocol       = "tcp"
  security_group_id = aws_security_group.alb.id
  from_port         = 80
  to_port           = 80
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_egress_rule" "lb_to_ecs_instance" {
  ip_protocol                  = "tcp"
  security_group_id            = aws_security_group.alb.id
  from_port                    = 8080
  to_port                      = 8080
  referenced_security_group_id = aws_security_group.ecs_instance.id
}

resource "aws_vpc_security_group_ingress_rule" "ecs_instance_from_alb" {
  ip_protocol                  = "tcp"
  security_group_id            = aws_security_group.ecs_instance.id
  from_port                    = 8080
  to_port                      = 8080
  referenced_security_group_id = aws_security_group.alb.id
}

resource "aws_vpc_security_group_egress_rule" "ecs_instance_to_https" {
  # ECSインスタンスから外部のHTTPSサーバーにアクセスするためのルール(ECR・SSM)
  ip_protocol       = "tcp"
  security_group_id = aws_security_group.ecs_instance.id
  from_port         = 443
  to_port           = 443
  cidr_ipv4         = "0.0.0.0/0"
}


# ALB周りのリソースを作成
resource "aws_lb" "alb" {
  name               = "${var.app_name}-alb"
  internal           = false
  load_balancer_type = "application" # ALB
  security_groups    = [aws_security_group.alb.id]
  subnets            = data.aws_subnets.public.ids
}

resource "aws_lb_target_group" "alb_target_group" {
  name        = "${var.app_name}-alb-target-group"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip" # IPアドレスをターゲットにする
  vpc_id      = data.aws_vpc.this.id
  health_check {
    healthy_threshold = 3
    interval          = 30
    path              = "/health_checks"
    protocol          = "HTTP"
    matcher           = "200"
    timeout           = 5
  }
}

resource "aws_lb_listener" "alb_listener" {
  load_balancer_arn = aws_lb.alb.arn
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Fixed response content"
      status_code  = "200"
    }
  }
}

resource "aws_lb_listener_rule" "alb_listener_rule" {
  listener_arn = aws_lb_listener.alb_listener.arn
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.alb_target_group.arn
  }
  condition {
    path_pattern {
      values = ["*"]
    }
  }
}


# ECSクラスターの作成
resource "aws_ecs_cluster" "ecs_cluster" {
  name = "${var.app_name}-ecs-cluster"
}

resource "aws_ecs_cluster_capacity_providers" "ecs_capacity_provider" {
  cluster_name       = aws_ecs_cluster.ecs_cluster.name
  capacity_providers = ["FARGATE"]
}


# タスク実行ロールを作成
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}


data "aws_ssm_parameters_by_path" "ssm_parameters" {
  # SSMパラメーターストアからパラメーターを取得
  path            = "/${var.app_name}/"
  recursive       = true
  with_decryption = true
}

data "aws_iam_policy_document" "ecs_task_execution_assume_role" {
  # ECSがRoleを引き受けるためのIAMポリシードキュメント
  statement {
    effect = "Allow"
    actions = [
      "sts:AssumeRole",
    ]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

data "aws_iam_policy" "managed_ecs_task_execution" {
  # ECSタスク実行用のIAMポリシー
  name = "AmazonECSTaskExecutionRolePolicy"
}

data "aws_iam_policy_document" "ecs_task_execution" {
  # ECSタスク実行用のIAMポリシードキュメント
  statement {
    effect = "Allow"
    actions = [
      "ssm:GetParameters",
      "ssm:GetParameter",
    ]
    resources = [
      for parameter in data.aws_ssm_parameters_by_path.ssm_parameters.names :
      "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter/${parameter}"
    ]
  }
}

resource "aws_iam_role" "ecs_task_execution_role" {
  # ECSタスク実行用のIAMロール
  name               = "${var.app_name}-ecs-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_execution_assume_role.json
}

resource "aws_iam_role_policy_attachments_exclusive" "ecs_task_execution_managed_policy" {
  # ECSタスク実行用のIAMロールにマネージドポリシーをアタッチ
  policy_arns = [data.aws_iam_policy.managed_ecs_task_execution.arn]
  role_name   = aws_iam_role.ecs_task_execution_role.name
}

resource "aws_iam_role_policy" "ecs_task_execution_inline_policy" {
  # ECSタスク実行用のIAMロールにインラインポリシーをアタッチ
  name   = "${var.app_name}-ecs-task-execution-inline-policy"
  role   = aws_iam_role.ecs_task_execution_role.name
  policy = data.aws_iam_policy_document.ecs_task_execution.json
}


# タスクロールを作成
data "aws_iam_policy_document" "ecs_task_assume_role" {
  statement {
    effect = "Allow"
    actions = [
      "sts:AssumeRole",
    ]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ecs_task" {
  # ECS Execの実行用インラインポリシー
  statement {
    effect = "Allow"
    actions = [
      "ssmmessages:CreateControlChannel",
      "ssmmessages:CreateDataChannel",
      "ssmmessages:OpenControlChannel",
      "ssmmessages:OpenDataChannel",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_role" "ecs_task" {
  name               = "${var.app_name}-ecs-task-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume_role.json
}

resource "aws_iam_role_policy" "ecs_task_inline_policy" {
  name   = "${var.app_name}-ecs-task-inline-policy"
  role   = aws_iam_role.ecs_task.name
  policy = data.aws_iam_policy_document.ecs_task.json
}


# ECSタスク定義を作成
data "aws_ecr_repository" "ecr_repository" {
  name = "${var.app_name}-ecr-repository"
}

resource "aws_cloudwatch_log_group" "ecs_log_group" {
  name              = "/ecs/${var.app_name}"
  retention_in_days = 30
}

resource "aws_ecs_task_definition" "ecs_task_definition" {
  family                   = var.app_name
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task.arn
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  requires_compatibilities = ["FARGATE"]
  skip_destroy             = true
  container_definitions = jsonencode([
    {
      name      = var.app_name
      image     = "medpeer/health_check:latest"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_log_group.name
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = var.app_name
        }

      }
      environment = [
        {
          name  = "NGINX_PORT"
          value = "8080"
        },
        {
          name  = "HEALTH_CHECK_PATH"
          value = "/health"
        }
      ]

    }
  ])
}

# ECSサービスを作成
resource "aws_ecs_service" "ecs_service" {
  name                              = "${var.app_name}-ecs-service"
  cluster                           = aws_ecs_cluster.ecs_cluster.id
  task_definition                   = aws_ecs_task_definition.ecs_task_definition.arn
  desired_count                     = 0
  enable_execute_command            = true # ECS Execを有効化
  health_check_grace_period_seconds = 60
  launch_type                       = "FARGATE"
  deployment_circuit_breaker {
    enable   = true
    rollback = false
  }
  network_configuration {
    subnets          = data.aws_subnets.private.ids
    security_groups  = [aws_security_group.ecs_instance.id]
    assign_public_ip = false
  }
  load_balancer {
    target_group_arn = aws_lb_target_group.alb_target_group.arn
    container_name   = var.app_name
    container_port   = 8080
  }
  lifecycle {
    ignore_changes = [desired_count] # desired_countの変更を無視
  }
  depends_on = [aws_lb_listener_rule.alb_listener_rule]
}
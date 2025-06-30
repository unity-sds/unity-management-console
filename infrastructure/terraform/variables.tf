# Variables for Unity Management Console ECS Fargate Deployment

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "unity"
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b"]
}

variable "desired_count" {
  description = "Desired number of ECS tasks"
  type        = number
  default     = 1
}

variable "max_capacity" {
  description = "Maximum number of ECS tasks for auto scaling"
  type        = number
  default     = 3
}

variable "cpu" {
  description = "CPU units for ECS task (256, 512, 1024, etc.)"
  type        = number
  default     = 512
}

variable "memory" {
  description = "Memory in MiB for ECS task"
  type        = number
  default     = 1024
}

variable "domain_name" {
  description = "Domain name for the application (optional)"
  type        = string
  default     = ""
}

variable "certificate_arn" {
  description = "ARN of the SSL certificate for HTTPS (optional)"
  type        = string
  default     = ""
}

variable "unity_config" {
  description = "Unity configuration parameters"
  type = object({
    awsregion           = string
    project             = string
    venue               = string
    installprefix       = string
    bucketname          = string
    marketplaceowner    = string
    marketplacerepo     = string
  })
  default = {
    awsregion           = "us-west-2"
    project             = "myproject"
    venue               = "dev"
    installprefix       = "unity"
    bucketname          = "unity-myproject-dev-terraform-state"
    marketplaceowner    = "unity-sds"
    marketplacerepo     = "unity-marketplace"
  }
}

variable "enable_container_insights" {
  description = "Enable CloudWatch Container Insights for ECS cluster"
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 30
}

variable "backup_schedule" {
  description = "Backup schedule for EFS (cron expression)"
  type        = string
  default     = "cron(0 2 * * ? *)" # Daily at 2 AM
}

variable "efs_throughput_mode" {
  description = "EFS throughput mode (provisioned or bursting)"
  type        = string
  default     = "provisioned"
}

variable "efs_provisioned_throughput" {
  description = "EFS provisioned throughput in MiB/s (only for provisioned mode)"
  type        = number
  default     = 10
}

variable "enable_efs_backup" {
  description = "Enable automatic EFS backups"
  type        = bool
  default     = true
}

variable "health_check_grace_period" {
  description = "Health check grace period in seconds"
  type        = number
  default     = 300
}

variable "deployment_configuration" {
  description = "ECS deployment configuration"
  type = object({
    maximum_percent                = number
    minimum_healthy_percent        = number
    deployment_circuit_breaker_enable = bool
    deployment_circuit_breaker_rollback = bool
  })
  default = {
    maximum_percent                = 200
    minimum_healthy_percent        = 100
    deployment_circuit_breaker_enable = true
    deployment_circuit_breaker_rollback = true
  }
}

variable "enable_execute_command" {
  description = "Enable ECS Exec for debugging"
  type        = bool
  default     = false
}

variable "platform_version" {
  description = "Fargate platform version"
  type        = string
  default     = "1.4.0"
}

variable "assign_public_ip" {
  description = "Assign public IP to ECS tasks (set to true if using public subnets)"
  type        = bool
  default     = false
}

variable "propagate_tags" {
  description = "How to propagate tags (TASK_DEFINITION or SERVICE)"
  type        = string
  default     = "SERVICE"
}

variable "enable_scheduling" {
  description = "Enable scheduled start/stop of ECS service"
  type        = bool
  default     = false
}

variable "schedule_start_time" {
  description = "Cron expression for starting the service (e.g., '0 8 * * MON-FRI')"
  type        = string
  default     = "0 8 * * MON-FRI" # 8 AM on weekdays
}

variable "schedule_stop_time" {
  description = "Cron expression for stopping the service (e.g., '0 18 * * MON-FRI')"
  type        = string
  default     = "0 18 * * MON-FRI" # 6 PM on weekdays
}

variable "ecr_scan_on_push" {
  description = "Enable ECR image scanning on push"
  type        = bool
  default     = true
}

variable "ecr_image_tag_mutability" {
  description = "ECR image tag mutability (MUTABLE or IMMUTABLE)"
  type        = string
  default     = "MUTABLE"
}

variable "additional_tags" {
  description = "Additional tags to apply to all resources"
  type        = map(string)
  default     = {}
}

variable "enable_deletion_protection" {
  description = "Enable deletion protection for ALB"
  type        = bool
  default     = false
}

variable "ssl_policy" {
  description = "SSL policy for HTTPS listener"
  type        = string
  default     = "ELBSecurityPolicy-TLS-1-2-2017-01"
}

variable "deregistration_delay" {
  description = "Time in seconds for ALB to wait before deregistering targets"
  type        = number
  default     = 300
}

variable "health_check_interval" {
  description = "Health check interval in seconds"
  type        = number
  default     = 30
}

variable "health_check_timeout" {
  description = "Health check timeout in seconds"
  type        = number
  default     = 5
}

variable "health_check_healthy_threshold" {
  description = "Number of consecutive successful health checks"
  type        = number
  default     = 2
}

variable "health_check_unhealthy_threshold" {
  description = "Number of consecutive failed health checks"
  type        = number
  default     = 2
}

variable "health_check_path" {
  description = "Health check path"
  type        = string
  default     = "/"
}

variable "health_check_matcher" {
  description = "HTTP codes to use when checking for a successful response"
  type        = string
  default     = "200"
}
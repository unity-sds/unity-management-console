# Outputs for Unity Management Console ECS Fargate Deployment

output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = aws_subnet.private[*].id
}

output "ecs_cluster_id" {
  description = "ID of the ECS cluster"
  value       = aws_ecs_cluster.main.id
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "ecs_service_name" {
  description = "Name of the ECS service"
  value       = aws_ecs_service.main.name
}

output "ecs_task_definition_arn" {
  description = "ARN of the ECS task definition"
  value       = aws_ecs_task_definition.main.arn
}

output "ecs_task_definition_family" {
  description = "Family of the ECS task definition"
  value       = aws_ecs_task_definition.main.family
}

output "ecs_task_definition_revision" {
  description = "Revision of the ECS task definition"
  value       = aws_ecs_task_definition.main.revision
}

output "load_balancer_dns_name" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "load_balancer_zone_id" {
  description = "Zone ID of the load balancer"
  value       = aws_lb.main.zone_id
}

output "load_balancer_arn" {
  description = "ARN of the load balancer"
  value       = aws_lb.main.arn
}

output "target_group_arn" {
  description = "ARN of the target group"
  value       = aws_lb_target_group.main.arn
}

output "ecr_repository_url" {
  description = "URL of the ECR repository"
  value       = aws_ecr_repository.main.repository_url
}

output "ecr_repository_name" {
  description = "Name of the ECR repository"
  value       = aws_ecr_repository.main.name
}

output "efs_file_system_id" {
  description = "ID of the EFS file system"
  value       = aws_efs_file_system.main.id
}

output "efs_file_system_dns_name" {
  description = "DNS name of the EFS file system"
  value       = aws_efs_file_system.main.dns_name
}

output "efs_mount_target_ids" {
  description = "IDs of the EFS mount targets"
  value       = aws_efs_mount_target.main[*].id
}

output "efs_mount_target_dns_names" {
  description = "DNS names of the EFS mount targets"
  value       = aws_efs_mount_target.main[*].dns_name
}

output "security_group_alb_id" {
  description = "ID of the ALB security group"
  value       = aws_security_group.alb.id
}

output "security_group_ecs_id" {
  description = "ID of the ECS security group"
  value       = aws_security_group.ecs.id
}

output "security_group_efs_id" {
  description = "ID of the EFS security group"
  value       = aws_security_group.efs.id
}

output "iam_role_ecs_task_execution_arn" {
  description = "ARN of the ECS task execution role"
  value       = aws_iam_role.ecs_task_execution.arn
}

output "iam_role_ecs_task_arn" {
  description = "ARN of the ECS task role"
  value       = aws_iam_role.ecs_task.arn
}

output "cloudwatch_log_group_name" {
  description = "Name of the CloudWatch log group"
  value       = aws_cloudwatch_log_group.main.name
}

output "cloudwatch_log_group_arn" {
  description = "ARN of the CloudWatch log group"
  value       = aws_cloudwatch_log_group.main.arn
}

output "nat_gateway_ids" {
  description = "IDs of the NAT gateways"
  value       = aws_nat_gateway.main[*].id
}

output "nat_gateway_public_ips" {
  description = "Public IPs of the NAT gateways"
  value       = aws_eip.nat[*].public_ip
}

output "application_url" {
  description = "URL to access the Unity Management Console"
  value       = var.certificate_arn != "" ? "https://${aws_lb.main.dns_name}" : "http://${aws_lb.main.dns_name}"
}

output "custom_domain_url" {
  description = "Custom domain URL (if configured)"
  value       = var.domain_name != "" ? "https://${var.domain_name}" : null
}

output "autoscaling_target_resource_id" {
  description = "Resource ID of the autoscaling target"
  value       = aws_appautoscaling_target.ecs_target.resource_id
}

output "deployment_commands" {
  description = "Commands to deploy the application"
  value = {
    build_and_push = "docker build -t ${aws_ecr_repository.main.repository_url}:latest . && docker push ${aws_ecr_repository.main.repository_url}:latest"
    update_service = "aws ecs update-service --cluster ${aws_ecs_cluster.main.name} --service ${aws_ecs_service.main.name} --force-new-deployment"
    start_service  = "aws ecs update-service --cluster ${aws_ecs_cluster.main.name} --service ${aws_ecs_service.main.name} --desired-count 1"
    stop_service   = "aws ecs update-service --cluster ${aws_ecs_cluster.main.name} --service ${aws_ecs_service.main.name} --desired-count 0"
  }
}

output "monitoring_urls" {
  description = "URLs for monitoring the deployment"
  value = {
    cloudwatch_logs = "https://${var.region}.console.aws.amazon.com/cloudwatch/home?region=${var.region}#logsV2:log-groups/log-group/${replace(aws_cloudwatch_log_group.main.name, "/", "%2F")}"
    ecs_cluster     = "https://${var.region}.console.aws.amazon.com/ecs/home?region=${var.region}#/clusters/${aws_ecs_cluster.main.name}"
    ecs_service     = "https://${var.region}.console.aws.amazon.com/ecs/home?region=${var.region}#/clusters/${aws_ecs_cluster.main.name}/services/${aws_ecs_service.main.name}"
    load_balancer   = "https://${var.region}.console.aws.amazon.com/ec2/v2/home?region=${var.region}#LoadBalancers:search=${aws_lb.main.name}"
    efs_file_system = "https://${var.region}.console.aws.amazon.com/efs/home?region=${var.region}#/file-systems/${aws_efs_file_system.main.id}"
  }
}

output "efs_access_point_commands" {
  description = "Commands to create EFS access points for specific directories"
  value = {
    create_workdir_access_point = "aws efs create-access-point --file-system-id ${aws_efs_file_system.main.id} --root-directory Path=/workdir --creation-info OwnerUid=1001,OwnerGid=1001,Permissions=755 --tags Key=Name,Value=workdir-access-point"
    create_database_access_point = "aws efs create-access-point --file-system-id ${aws_efs_file_system.main.id} --root-directory Path=/database --creation-info OwnerUid=1001,OwnerGid=1001,Permissions=755 --tags Key=Name,Value=database-access-point"
    create_config_access_point = "aws efs create-access-point --file-system-id ${aws_efs_file_system.main.id} --root-directory Path=/config --creation-info OwnerUid=1001,OwnerGid=1001,Permissions=755 --tags Key=Name,Value=config-access-point"
  }
}

output "backup_commands" {
  description = "Commands for backup operations"
  value = {
    create_manual_backup = "aws efs create-backup --file-system-id ${aws_efs_file_system.main.id} --creation-token unity-backup-$(date +%Y%m%d_%H%M%S)"
    list_backups = "aws efs describe-backups --file-system-id ${aws_efs_file_system.main.id}"
  }
}

output "troubleshooting_commands" {
  description = "Commands for troubleshooting"
  value = {
    check_service_status = "aws ecs describe-services --cluster ${aws_ecs_cluster.main.name} --services ${aws_ecs_service.main.name}"
    check_task_status = "aws ecs list-tasks --cluster ${aws_ecs_cluster.main.name} --service-name ${aws_ecs_service.main.name}"
    check_efs_mount_targets = "aws efs describe-mount-targets --file-system-id ${aws_efs_file_system.main.id}"
    check_target_health = "aws elbv2 describe-target-health --target-group-arn ${aws_lb_target_group.main.arn}"
    view_logs = "aws logs tail ${aws_cloudwatch_log_group.main.name} --follow"
  }
}
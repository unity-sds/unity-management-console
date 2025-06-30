# Unity Management Console ECS Fargate Deployment Guide

This guide provides step-by-step instructions for deploying the Unity Management Console on AWS ECS Fargate with EFS storage. This deployment model allows for on-demand scaling and persistent storage, making it ideal for environments where the console doesn't need to run continuously.

## Architecture Overview

The ECS Fargate deployment includes:
- **ECS Fargate Service**: Runs the containerized Unity Management Console
- **EFS File System**: Provides persistent storage for database, workdir, and configuration
- **Application Load Balancer**: Provides HTTPS termination and routing
- **VPC with Public/Private Subnets**: Network isolation and security
- **EFS Mount Targets**: Enable EFS access from Fargate tasks
- **Security Groups**: Control network access between components

## Benefits of ECS Fargate Deployment

- **On-demand scaling**: Start/stop the service as needed
- **Persistent storage**: Database and workdir survive container restarts
- **No server management**: Fully managed container platform
- **Cost-effective**: Pay only when running
- **High availability**: Built-in redundancy across AZs
- **Integrated monitoring**: CloudWatch logs and metrics

## Prerequisites

- AWS Account with appropriate permissions
- AWS CLI configured with credentials
- Docker installed (for building custom images)
- Terraform installed (for infrastructure deployment)
- Domain name for HTTPS access (optional but recommended)

## Deployment Options

This guide provides two deployment methods:

1. **Quick Deploy**: Use pre-built infrastructure templates
2. **Custom Deploy**: Build and customize your own infrastructure

## Option 1: Quick Deploy with Terraform

### 1. Clone the Repository

```bash
git clone https://github.com/unity-sds/unity-management-console.git
cd unity-management-console
```

### 2. Configure Deployment Variables

Create a `terraform.tfvars` file:

```hcl
# Basic Configuration
project_name = "unity"
environment  = "dev"
region      = "us-west-2"

# Networking
vpc_cidr = "10.0.0.0/16"
availability_zones = ["us-west-2a", "us-west-2b"]

# ECS Configuration
desired_count = 1
cpu          = 512
memory       = 1024

# Domain Configuration (optional)
domain_name = "unity-console.yourdomain.com"
certificate_arn = "arn:aws:acm:us-west-2:123456789012:certificate/12345678-1234-1234-1234-123456789012"

# Unity Configuration
unity_config = {
  awsregion           = "us-west-2"
  project             = "myproject"
  venue               = "dev"
  installprefix       = "unity"
  bucketname          = "unity-myproject-dev-terraform-state"
  marketplaceowner    = "unity-sds"
  marketplacerepo     = "unity-marketplace"
}
```

### 3. Deploy Infrastructure

```bash
cd infrastructure/terraform
terraform init
terraform plan
terraform apply
```

### 4. Access the Console

After deployment, the console will be available at:
- **With custom domain**: `https://unity-console.yourdomain.com`
- **With ALB DNS**: `https://<alb-dns-name>.us-west-2.elb.amazonaws.com`

## Option 2: Custom Deploy

### 1. Build the Docker Image

Create a Dockerfile in the project root:

```dockerfile
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN go build -o management-console cmd/web/main.go

FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY ui/ .
RUN npm install
RUN npm run build

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Install required tools
RUN apk add --no-cache \
    curl \
    wget \
    unzip \
    git \
    bash

# Install Terraform
RUN wget https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip && \
    unzip terraform_1.5.7_linux_amd64.zip && \
    mv terraform /usr/local/bin/ && \
    rm terraform_1.5.7_linux_amd64.zip

# Install AWS CLI
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install && \
    rm -rf awscliv2.zip aws

# Copy built applications
COPY --from=backend-builder /app/management-console .
COPY --from=frontend-builder /app/build ./ui/build

# Create directories
RUN mkdir -p /data/workdir
RUN mkdir -p /data/database

# Expose port
EXPOSE 8080

# Set environment variables
ENV UNITY_WORKDIR=/data/workdir
ENV UNITY_DATABASE_PATH=/data/database

CMD ["./management-console", "webapp"]
```

### 2. Build and Push Image

```bash
# Build the image
docker build -t unity-management-console .

# Tag for ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-west-2.amazonaws.com

# Create ECR repository
aws ecr create-repository --repository-name unity-management-console --region us-west-2

# Tag and push
docker tag unity-management-console:latest <account-id>.dkr.ecr.us-west-2.amazonaws.com/unity-management-console:latest
docker push <account-id>.dkr.ecr.us-west-2.amazonaws.com/unity-management-console:latest
```

### 3. Create EFS File System

```bash
# Create EFS file system
aws efs create-file-system \
    --creation-token unity-console-efs \
    --performance-mode generalPurpose \
    --throughput-mode provisioned \
    --provisioned-throughput-in-mibps 10 \
    --tags Key=Name,Value=unity-console-efs
```

### 4. Create ECS Task Definition

Save as `task-definition.json`:

```json
{
  "family": "unity-management-console",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::<account-id>:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::<account-id>:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "unity-console",
      "image": "<account-id>.dkr.ecr.us-west-2.amazonaws.com/unity-management-console:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "UNITY_WORKDIR",
          "value": "/data/workdir"
        },
        {
          "name": "UNITY_DATABASE_PATH",
          "value": "/data/database"
        },
        {
          "name": "UNITY_AWSREGION",
          "value": "us-west-2"
        }
      ],
      "mountPoints": [
        {
          "sourceVolume": "efs-storage",
          "containerPath": "/data",
          "readOnly": false
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/unity-management-console",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ],
  "volumes": [
    {
      "name": "efs-storage",
      "efsVolumeConfiguration": {
        "fileSystemId": "<efs-file-system-id>",
        "rootDirectory": "/",
        "transitEncryption": "ENABLED"
      }
    }
  ]
}
```

## Configuration Management

### EFS Directory Structure

The EFS file system will contain:

```
/data/
├── workdir/           # Unity working directory
│   ├── workspace/     # Terraform workspaces
│   ├── install_logs/  # Installation logs
│   └── module_cache/  # Terraform modules
├── database/          # Database files
│   └── test.db        # SQLite database
└── config/            # Configuration files
    └── unity.yaml     # Unity configuration
```

### Unity Configuration

Create the Unity configuration file that will be mounted into the container:

```yaml
# /data/config/unity.yaml
awsregion: "us-west-2"
bucketname: "unity-<project>-<venue>-terraform-state"
project: "myproject"
venue: "dev"
installprefix: "unity"
consolehost: "https://unity-console.yourdomain.com"
basepath: "https://unity-console.yourdomain.com"
workdir: "/data/workdir"
marketplacebaseurl: "https://raw.githubusercontent.com/"
marketplaceowner: "unity-sds"
marketplacerepo: "unity-marketplace"

MarketplaceItems:
  - name: unity-proxy
    version: "0.5.0"
  - name: unity-apigateway
    version: "0.3.0"
  - name: unity-cs-monitoring-lambda
    version: "1.1.0"
  - name: unity-portal
    version: "0.1.0"

defaultssmparameters:
  - name: "/unity/project"
    value: "myproject"
  - name: "/unity/venue"
    value: "dev"
```

## Infrastructure Components

### VPC and Networking

```bash
# Create VPC
aws ec2 create-vpc --cidr-block 10.0.0.0/16 --tag-specifications 'ResourceType=vpc,Tags=[{Key=Name,Value=unity-console-vpc}]'

# Create subnets
aws ec2 create-subnet --vpc-id <vpc-id> --cidr-block 10.0.1.0/24 --availability-zone us-west-2a
aws ec2 create-subnet --vpc-id <vpc-id> --cidr-block 10.0.2.0/24 --availability-zone us-west-2b
```

### Security Groups

```bash
# ECS Security Group
aws ec2 create-security-group \
    --group-name unity-console-ecs \
    --description "Unity Console ECS Security Group" \
    --vpc-id <vpc-id>

# Allow inbound from ALB
aws ec2 authorize-security-group-ingress \
    --group-id <ecs-sg-id> \
    --protocol tcp \
    --port 8080 \
    --source-group <alb-sg-id>

# EFS Security Group
aws ec2 create-security-group \
    --group-name unity-console-efs \
    --description "Unity Console EFS Security Group" \
    --vpc-id <vpc-id>

# Allow NFS from ECS
aws ec2 authorize-security-group-ingress \
    --group-id <efs-sg-id> \
    --protocol tcp \
    --port 2049 \
    --source-group <ecs-sg-id>
```

### IAM Roles

The ECS tasks require specific IAM roles:

**ECS Task Execution Role**:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    }
  ]
}
```

**ECS Task Role** (for Unity operations):
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:*",
        "dynamodb:*",
        "ssm:*",
        "iam:*",
        "lambda:*",
        "apigateway:*",
        "cloudformation:*",
        "ec2:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Monitoring and Logging

### CloudWatch Logs

```bash
# Create log group
aws logs create-log-group --log-group-name /ecs/unity-management-console
```

### CloudWatch Metrics

Key metrics to monitor:
- CPU utilization
- Memory utilization
- EFS throughput
- HTTP response times
- Error rates

## Scaling and Management

### Auto Scaling

```json
{
  "ServiceName": "unity-management-console",
  "ClusterName": "unity-cluster",
  "ScalableDimension": "ecs:service:DesiredCount",
  "ServiceNamespace": "ecs",
  "MinCapacity": 0,
  "MaxCapacity": 3
}
```

### Start/Stop Operations

```bash
# Start the service
aws ecs update-service \
    --cluster unity-cluster \
    --service unity-management-console \
    --desired-count 1

# Stop the service
aws ecs update-service \
    --cluster unity-cluster \
    --service unity-management-console \
    --desired-count 0
```

## Backup and Recovery

### EFS Backup

```bash
# Enable automatic backups
aws efs put-backup-policy \
    --file-system-id <efs-id> \
    --backup-policy Status=ENABLED
```

### Database Backup

```bash
# Create backup script
#!/bin/bash
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
aws efs create-backup \
    --file-system-id <efs-id> \
    --creation-token unity-backup-${TIMESTAMP}
```

## Security Considerations

1. **Network Security**:
   - Use private subnets for ECS tasks
   - Implement security groups with least privilege
   - Enable VPC Flow Logs

2. **Data Encryption**:
   - Enable EFS encryption at rest
   - Use HTTPS for all communications
   - Encrypt ECS task communication

3. **IAM Security**:
   - Use least privilege IAM roles
   - Regularly audit permissions
   - Enable CloudTrail logging

4. **Container Security**:
   - Use minimal base images
   - Regularly update dependencies
   - Scan images for vulnerabilities

## Troubleshooting

### Common Issues

1. **EFS Mount Failures**:
   ```bash
   # Check EFS mount targets
   aws efs describe-mount-targets --file-system-id <efs-id>
   
   # Verify security group rules
   aws ec2 describe-security-groups --group-ids <efs-sg-id>
   ```

2. **Task Startup Failures**:
   ```bash
   # Check ECS service events
   aws ecs describe-services --cluster unity-cluster --services unity-management-console
   
   # View container logs
   aws logs get-log-events --log-group-name /ecs/unity-management-console --log-stream-name <stream-name>
   ```

3. **Database Connection Issues**:
   ```bash
   # Verify EFS permissions
   aws efs describe-file-systems --file-system-id <efs-id>
   
   # Check mount point permissions
   docker exec -it <container-id> ls -la /data/database/
   ```

### Health Checks

```bash
# ECS service health
aws ecs describe-services --cluster unity-cluster --services unity-management-console

# ALB target health
aws elbv2 describe-target-health --target-group-arn <target-group-arn>

# EFS health
aws efs describe-file-systems --file-system-id <efs-id>
```

## Cost Optimization

1. **Fargate Savings Plans**: Use Compute Savings Plans for predictable workloads
2. **Spot Capacity**: Consider Fargate Spot for non-critical environments
3. **EFS Intelligent Tiering**: Enable to automatically move data to cheaper storage classes
4. **Scheduled Scaling**: Automatically scale down during off-hours

## Migration from EC2

To migrate from an existing EC2 deployment:

1. **Backup existing data**:
   ```bash
   # Backup database
   cp /home/ubuntu/unity-management-console/backend/test.db ~/backup/
   
   # Backup configuration
   cp ~/.unity/unity.yaml ~/backup/
   
   # Backup workdir
   tar -czf ~/backup/workdir.tar.gz /home/ubuntu/unity-workdir/
   ```

2. **Upload to EFS**:
   ```bash
   # Mount EFS locally
   sudo mount -t efs <efs-id>:/ /mnt/efs
   
   # Copy data
   sudo cp ~/backup/test.db /mnt/efs/database/
   sudo cp ~/backup/unity.yaml /mnt/efs/config/
   sudo tar -xzf ~/backup/workdir.tar.gz -C /mnt/efs/
   ```

3. **Deploy ECS service** using the steps above

## Additional Resources

- [AWS ECS Documentation](https://docs.aws.amazon.com/ecs/)
- [AWS EFS Documentation](https://docs.aws.amazon.com/efs/)
- [Unity SDS Documentation](https://unity-sds.gitbook.io/docs/)
- [Unity Marketplace](https://github.com/unity-sds/unity-marketplace)
- [Management Console Repository](https://github.com/unity-sds/unity-management-console)
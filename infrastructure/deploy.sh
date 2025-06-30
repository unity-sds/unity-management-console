#!/bin/bash

# Unity Management Console ECS Fargate Deployment Script
# This script automates the deployment of the Unity Management Console to ECS Fargate

set -e

# Configuration
PROJECT_NAME=${PROJECT_NAME:-"unity"}
ENVIRONMENT=${ENVIRONMENT:-"dev"}
AWS_REGION=${AWS_REGION:-"us-west-2"}
DOCKER_TAG=${DOCKER_TAG:-"latest"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check AWS CLI
    if ! command -v aws &> /dev/null; then
        log_error "AWS CLI is not installed. Please install it first."
        exit 1
    fi
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install it first."
        exit 1
    fi
    
    # Check Terraform (if using Terraform deployment)
    if [ "$DEPLOYMENT_METHOD" = "terraform" ] && ! command -v terraform &> /dev/null; then
        log_error "Terraform is not installed. Please install it first."
        exit 1
    fi
    
    # Check AWS credentials
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS credentials are not configured. Please run 'aws configure' first."
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Get AWS account ID
get_account_id() {
    aws sts get-caller-identity --query Account --output text
}

# Create ECR repository if it doesn't exist
create_ecr_repository() {
    local repo_name="$1"
    
    log_info "Checking if ECR repository exists: $repo_name"
    
    if ! aws ecr describe-repositories --repository-names "$repo_name" --region "$AWS_REGION" &> /dev/null; then
        log_info "Creating ECR repository: $repo_name"
        aws ecr create-repository \
            --repository-name "$repo_name" \
            --image-scanning-configuration scanOnPush=true \
            --region "$AWS_REGION"
        log_success "ECR repository created: $repo_name"
    else
        log_info "ECR repository already exists: $repo_name"
    fi
}

# Build and push Docker image
build_and_push_image() {
    local account_id="$1"
    local repo_name="$2"
    local tag="$3"
    
    local image_uri="${account_id}.dkr.ecr.${AWS_REGION}.amazonaws.com/${repo_name}:${tag}"
    
    log_info "Building Docker image..."
    docker build -t "$image_uri" .
    
    log_info "Logging in to ECR..."
    aws ecr get-login-password --region "$AWS_REGION" | \
        docker login --username AWS --password-stdin "${account_id}.dkr.ecr.${AWS_REGION}.amazonaws.com"
    
    log_info "Pushing image to ECR: $image_uri"
    docker push "$image_uri"
    
    log_success "Image pushed successfully: $image_uri"
    echo "$image_uri"
}

# Deploy infrastructure using Terraform
deploy_terraform() {
    local tfvars_file="$1"
    
    log_info "Deploying infrastructure using Terraform..."
    
    cd infrastructure/terraform
    
    # Initialize Terraform
    log_info "Initializing Terraform..."
    terraform init
    
    # Plan deployment
    log_info "Planning Terraform deployment..."
    terraform plan -var-file="$tfvars_file"
    
    # Apply deployment
    log_info "Applying Terraform deployment..."
    terraform apply -var-file="$tfvars_file" -auto-approve
    
    cd ../..
    
    log_success "Infrastructure deployed successfully"
}

# Deploy infrastructure using CloudFormation
deploy_cloudformation() {
    local stack_name="${PROJECT_NAME}-${ENVIRONMENT}-unity-console"
    local template_file="infrastructure/cloudformation/unity-console-infrastructure.yaml"
    local parameters_file="$1"
    
    log_info "Deploying infrastructure using CloudFormation..."
    
    # Check if stack exists
    if aws cloudformation describe-stacks --stack-name "$stack_name" --region "$AWS_REGION" &> /dev/null; then
        log_info "Updating existing CloudFormation stack: $stack_name"
        aws cloudformation update-stack \
            --stack-name "$stack_name" \
            --template-body "file://$template_file" \
            --parameters "file://$parameters_file" \
            --capabilities CAPABILITY_NAMED_IAM \
            --region "$AWS_REGION"
    else
        log_info "Creating new CloudFormation stack: $stack_name"
        aws cloudformation create-stack \
            --stack-name "$stack_name" \
            --template-body "file://$template_file" \
            --parameters "file://$parameters_file" \
            --capabilities CAPABILITY_NAMED_IAM \
            --region "$AWS_REGION"
    fi
    
    log_info "Waiting for CloudFormation stack to complete..."
    aws cloudformation wait stack-create-complete \
        --stack-name "$stack_name" \
        --region "$AWS_REGION" || \
    aws cloudformation wait stack-update-complete \
        --stack-name "$stack_name" \
        --region "$AWS_REGION"
    
    log_success "Infrastructure deployed successfully"
}

# Update ECS service
update_ecs_service() {
    local cluster_name="${PROJECT_NAME}-${ENVIRONMENT}-cluster"
    local service_name="${PROJECT_NAME}-${ENVIRONMENT}-console"
    
    log_info "Updating ECS service..."
    
    aws ecs update-service \
        --cluster "$cluster_name" \
        --service "$service_name" \
        --force-new-deployment \
        --region "$AWS_REGION"
    
    log_info "Waiting for service to stabilize..."
    aws ecs wait services-stable \
        --cluster "$cluster_name" \
        --services "$service_name" \
        --region "$AWS_REGION"
    
    log_success "ECS service updated successfully"
}

# Start ECS service
start_service() {
    local cluster_name="${PROJECT_NAME}-${ENVIRONMENT}-cluster"
    local service_name="${PROJECT_NAME}-${ENVIRONMENT}-console"
    local desired_count="${DESIRED_COUNT:-1}"
    
    log_info "Starting ECS service with desired count: $desired_count"
    
    aws ecs update-service \
        --cluster "$cluster_name" \
        --service "$service_name" \
        --desired-count "$desired_count" \
        --region "$AWS_REGION"
    
    log_success "ECS service start command sent"
}

# Stop ECS service
stop_service() {
    local cluster_name="${PROJECT_NAME}-${ENVIRONMENT}-cluster"
    local service_name="${PROJECT_NAME}-${ENVIRONMENT}-console"
    
    log_info "Stopping ECS service..."
    
    aws ecs update-service \
        --cluster "$cluster_name" \
        --service "$service_name" \
        --desired-count 0 \
        --region "$AWS_REGION"
    
    log_success "ECS service stop command sent"
}

# Get service status
get_service_status() {
    local cluster_name="${PROJECT_NAME}-${ENVIRONMENT}-cluster"
    local service_name="${PROJECT_NAME}-${ENVIRONMENT}-console"
    
    log_info "Getting service status..."
    
    aws ecs describe-services \
        --cluster "$cluster_name" \
        --services "$service_name" \
        --region "$AWS_REGION" \
        --query 'services[0].{Status:status,Running:runningCount,Desired:desiredCount,Pending:pendingCount}' \
        --output table
}

# Get application URL
get_application_url() {
    local stack_name="${PROJECT_NAME}-${ENVIRONMENT}-unity-console"
    
    log_info "Getting application URL..."
    
    if [ "$DEPLOYMENT_METHOD" = "cloudformation" ]; then
        aws cloudformation describe-stacks \
            --stack-name "$stack_name" \
            --region "$AWS_REGION" \
            --query 'Stacks[0].Outputs[?OutputKey==`ApplicationURL`].OutputValue' \
            --output text
    else
        # For Terraform, get ALB DNS name
        cd infrastructure/terraform
        terraform output -raw application_url
        cd ../..
    fi
}

# Show help
show_help() {
    cat << EOF
Unity Management Console ECS Fargate Deployment Script

Usage: $0 [COMMAND] [OPTIONS]

Commands:
    deploy              Deploy the complete infrastructure and application
    build               Build and push Docker image only
    update              Update the ECS service with latest image
    start               Start the ECS service
    stop                Stop the ECS service
    status              Show service status
    url                 Get application URL
    help                Show this help message

Options:
    --method <terraform|cloudformation>  Deployment method (default: terraform)
    --config <file>                      Configuration file
    --project <name>                     Project name (default: unity)
    --environment <env>                  Environment (default: dev)
    --region <region>                    AWS region (default: us-west-2)
    --tag <tag>                          Docker image tag (default: latest)

Environment Variables:
    PROJECT_NAME        Project name
    ENVIRONMENT         Environment name
    AWS_REGION          AWS region
    DOCKER_TAG          Docker image tag
    DEPLOYMENT_METHOD   Deployment method (terraform or cloudformation)
    DESIRED_COUNT       Desired number of tasks when starting service

Examples:
    $0 deploy --method terraform --config terraform.tfvars
    $0 deploy --method cloudformation --config parameters.json
    $0 build --tag v1.0.0
    $0 update
    $0 start
    $0 stop
    $0 status
    $0 url

EOF
}

# Main function
main() {
    local command="$1"
    shift
    
    # Parse options
    while [[ $# -gt 0 ]]; do
        case $1 in
            --method)
                DEPLOYMENT_METHOD="$2"
                shift 2
                ;;
            --config)
                CONFIG_FILE="$2"
                shift 2
                ;;
            --project)
                PROJECT_NAME="$2"
                shift 2
                ;;
            --environment)
                ENVIRONMENT="$2"
                shift 2
                ;;
            --region)
                AWS_REGION="$2"
                shift 2
                ;;
            --tag)
                DOCKER_TAG="$2"
                shift 2
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Set defaults
    DEPLOYMENT_METHOD=${DEPLOYMENT_METHOD:-"terraform"}
    ACCOUNT_ID=$(get_account_id)
    REPO_NAME="${PROJECT_NAME}-${ENVIRONMENT}-console"
    
    case $command in
        deploy)
            check_prerequisites
            create_ecr_repository "$REPO_NAME"
            build_and_push_image "$ACCOUNT_ID" "$REPO_NAME" "$DOCKER_TAG"
            
            if [ "$DEPLOYMENT_METHOD" = "terraform" ]; then
                deploy_terraform "${CONFIG_FILE:-terraform.tfvars}"
            elif [ "$DEPLOYMENT_METHOD" = "cloudformation" ]; then
                deploy_cloudformation "${CONFIG_FILE:-parameters.json}"
            else
                log_error "Invalid deployment method: $DEPLOYMENT_METHOD"
                exit 1
            fi
            
            log_success "Deployment completed successfully!"
            log_info "Application URL: $(get_application_url)"
            ;;
        build)
            check_prerequisites
            create_ecr_repository "$REPO_NAME"
            build_and_push_image "$ACCOUNT_ID" "$REPO_NAME" "$DOCKER_TAG"
            ;;
        update)
            update_ecs_service
            ;;
        start)
            start_service
            ;;
        stop)
            stop_service
            ;;
        status)
            get_service_status
            ;;
        url)
            get_application_url
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
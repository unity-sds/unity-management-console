# Unity Management Console EC2 Deployment Guide

This guide provides step-by-step instructions for manually deploying the Unity Management Console on an AWS EC2 instance.

## Installation Directory Structure

Throughout this guide, we use the following placeholder paths:
- `<INSTALL_DIR>`: The directory where Unity Management Console is installed (default: `/home/ubuntu/unity-management-console`)
- `<WORKDIR>`: The working directory for Unity operations (configurable via `~/.unity/unity.yaml`, default: `/home/ubuntu/unity-workdir`)

**Important Path Behavior Notes:**
- **Configurable paths**: Working directory (`<WORKDIR>`), user config directory (`~/.unity`)
- **Fixed relative paths**: Database file (`test.db`) is created in the directory where the service starts (typically `<INSTALL_DIR>`)
- **Hardcoded system paths**: The systemd service file contains absolute paths that must match your actual installation location

## Prerequisites

- AWS Account with appropriate permissions
- EC2 instance running Ubuntu 20.04 or later
- SSH access to the EC2 instance
- AWS CLI configured with credentials
- Minimum instance size: t3.medium (2 vCPU, 4GB RAM)
- Security group with ports 8080 (console) and 80/443 (if using proxy) open

## Installation Steps

### 1. Connect to Your EC2 Instance

```bash
ssh -i your-key.pem ubuntu@your-ec2-instance-ip
```

### 2. Install System Dependencies

```bash
# Update system packages
sudo apt-get update
sudo apt-get upgrade -y

# Install required dependencies
sudo apt-get install -y \
    build-essential \
    git \
    curl \
    wget \
    unzip \
    python3-pip \
    nodejs \
    npm

# Install Go (required for building from source)
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/bin/go' >> ~/.bashrc
source ~/.bashrc

# Install Terraform
wget https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip
unzip terraform_1.5.7_linux_amd64.zip
sudo mv terraform /usr/local/bin/
sudo chmod +x /usr/local/bin/terraform
```

### 3. Install AWS CLI

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

### 4. Configure AWS Credentials

```bash
aws configure
# Enter your AWS Access Key ID
# Enter your AWS Secret Access Key
# Enter your default region (e.g., us-west-2)
# Enter your default output format (json)
```

### 5. Clone and Build the Management Console

```bash
# Clone the repository
git clone https://github.com/unity-sds/unity-management-console.git
cd unity-management-console

# Build all
npm run build-all

# Build the backend
cd backend
go build -o management-console cmd/web/main.go
cd ..

# Build the frontend
cd ui
npm install
npm run build
cd ..
```

### 6. Create Unity Configuration Directory

```bash
mkdir -p ~/.unity
```

### 7. Configure the Management Console

Create the configuration file `~/.unity/unity.yaml`:

```yaml
# AWS Configuration
awsregion: "us-west-2"
bucketname: "unity-<project>-<venue>-terraform-state"

# Project Configuration
project: "myproject"
venue: "dev"
installprefix: "unity"

# Console Configuration
consolehost: "http://your-ec2-instance-ip:8080"
basepath: "http://your-ec2-instance-ip"

# Working Directory
workdir: "<WORKDIR>"

# Marketplace Configuration
marketplacebaseurl: "https://raw.githubusercontent.com/"
marketplaceowner: "unity-sds"
marketplacerepo: "unity-marketplace"

# Core Applications for Bootstrap
MarketplaceItems:
  - name: unity-proxy
    version: "0.5.0"
  - name: unity-apigateway
    version: "0.3.0"
  - name: unity-cs-monitoring-lambda
    version: "1.1.0"
  - name: unity-portal
    version: "0.1.0"

# SSM Parameters (optional)
defaultssmparameters:
  - name: "/unity/project"
    value: "myproject"
  - name: "/unity/venue"
    value: "dev"
```

### 8. Create Working Directory

```bash
mkdir -p <WORKDIR>
```

### 9. Set Up as a System Service

Create a systemd service file:

```bash
sudo tee /etc/systemd/system/unity-management-console.service << EOF
[Unit]
Description=Unity Management Console
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=<INSTALL_DIR>
ExecStart=<INSTALL_DIR>/backend/management-console webapp
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=unity-console

# Environment variables
Environment="PATH=/usr/local/bin:/usr/bin:/bin"
Environment="HOME=/home/ubuntu"

[Install]
WantedBy=multi-user.target
EOF
```

**Important**: Replace `<INSTALL_DIR>` with your actual installation path (e.g., `/home/ubuntu/unity-management-console`). The systemd service file requires absolute paths and cannot use placeholder variables.

### 10. Start the Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable the service to start on boot
sudo systemctl enable unity-management-console

# Start the service
sudo systemctl start unity-management-console

# Check service status
sudo systemctl status unity-management-console
```

### 11. Bootstrap the Environment

On first run, bootstrap the environment:

```bash
cd <INSTALL_DIR>/backend
./management-console webapp --bootstrap
```

This will:
- Initialize the local database
- Create S3 bucket for Terraform state
- Create DynamoDB table for state locking
- Set up SSM parameters
- Install core Unity applications

### 12. Access the Console

Open your web browser and navigate to:
```
http://your-ec2-instance-ip:8080
```

## Configuration Options

### Marketplace Configuration

The marketplace configuration determines where application packages are sourced from:

- **marketplacebaseurl**: Base URL for the marketplace (default: GitHub raw content)
- **marketplaceowner**: GitHub organization/user owning the marketplace repository
- **marketplacerepo**: Name of the marketplace repository

To use a custom marketplace:

1. Fork the unity-marketplace repository
2. Update the configuration:
   ```yaml
   marketplaceowner: "your-github-org"
   marketplacerepo: "your-marketplace-repo"
   ```

### Bootstrap Configuration

The bootstrap process installs core applications defined in `MarketplaceItems`:

```yaml
MarketplaceItems:
  - name: application-name
    version: "x.x.x"
```

Add or modify entries to customize which applications are installed during bootstrap.

### Module Registry Configuration

The Unity Management Console supports a Module Registry that provides reusable Terraform modules for applications. To enable this feature:

#### Option 1: Local Module Registry

Create a `module-registry.json` file in your workdir:

```bash
sudo -u ubuntu tee <WORKDIR>/module-registry.json << 'EOF'
{
  "version": "1.0",
  "metadata": {
    "name": "Unity Terraform Modules",
    "description": "Reusable Terraform modules for Unity applications",
    "last_updated": "2025-06-30T10:00:00Z"
  },
  "modules": {
    "unity-vpc": {
      "description": "Standard Unity VPC with public/private subnets",
      "source": "github.com/unity-sds/terraform-modules//modules/vpc",
      "versions": {
        "1.0.0": { "ref": "v1.0.0" },
        "latest": { "ref": "main" }
      }
    }
  }
}
EOF
```

#### Option 2: Marketplace Module Registry

The module registry will be automatically downloaded from your marketplace repository:

```yaml
# In your unity.yaml configuration
marketplaceowner: "unity-sds"
marketplacerepo: "unity-marketplace"
```

The console will look for:
```
https://raw.githubusercontent.com/unity-sds/unity-marketplace/main/module-registry.json
```

#### Using Modules in Applications

Applications can reference modules during installation by including a `moduleReferences` section in their installation parameters. See the [Module Registry Guide](module-registry-guide.md) for detailed usage instructions.

### AWS Configuration

- **awsregion**: AWS region for resources
- **bucketname**: S3 bucket for Terraform state (created during bootstrap)
- **project**: Unity project name
- **venue**: Deployment environment (dev, test, prod)

### Advanced Configuration

#### Using Environment Variables

All configuration options can be set via environment variables:

```bash
export UNITY_AWSREGION="us-east-1"
export UNITY_PROJECT="myproject"
export UNITY_VENUE="prod"
```

#### Custom Working Directory

```yaml
workdir: "/data/unity"
```

Ensure the directory exists and has appropriate permissions.

## Monitoring and Logs

### View Application Logs

```bash
# Real-time logs
sudo journalctl -u unity-management-console -f

# Last 100 lines
sudo journalctl -u unity-management-console -n 100
```

### Installation Logs

Installation logs are stored in:
```
<workdir>/install_logs/
```

### Database Location

The SQLite database is created as `test.db` in the working directory where the service starts:
```
<INSTALL_DIR>/backend/test.db
```

**Note**: The database filename `test.db` is hardcoded in the application. If you change the service's working directory, the database location will change accordingly.

## Security Considerations

1. **IAM Permissions**: Ensure the EC2 instance has an IAM role with necessary permissions:
   - S3 access for Terraform state
   - DynamoDB access for state locking
   - SSM Parameter Store access
   - Permissions to create resources for installed applications

2. **Security Groups**: Configure security groups to restrict access:
   ```bash
   # Example security group rules
   - Port 8080: Your IP or VPN CIDR only
   - Port 22: Your IP for SSH access
   - Port 80/443: If using unity-proxy
   ```

3. **HTTPS Setup**: For production, use a reverse proxy (nginx) with SSL:
   ```nginx
   server {
       listen 443 ssl;
       server_name your-domain.com;
       
       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

## Troubleshooting

### Service Won't Start

1. Check logs:
   ```bash
   sudo journalctl -u unity-management-console -n 50
   ```

2. Verify configuration:
   ```bash
   cat ~/.unity/unity.yaml
   ```

3. Check permissions:
   ```bash
   ls -la <WORKDIR>
   ```

### Bootstrap Fails

1. Verify AWS credentials:
   ```bash
   aws sts get-caller-identity
   ```

2. Check IAM permissions for S3 and DynamoDB

3. Ensure MarketplaceItems are properly configured

### Module Registry Issues

1. Registry not loading:
   ```bash
   # Check if registry file exists
   ls -la <WORKDIR>/module-registry.json
   
   # Check logs for registry loading errors
   sudo journalctl -u unity-management-console | grep -i "registry"
   ```

2. Module download fails:
   ```bash
   # Verify git is installed
   which git
   
   # Check module cache directory
   ls -la <WORKDIR>/module_cache/
   ```

3. Module not found:
   - Verify module name matches exactly in registry
   - Check requested version exists
   - Ensure source repository is accessible

### Cannot Access UI

1. Check security group rules
2. Verify service is running:
   ```bash
   sudo systemctl status unity-management-console
   ```
3. Test locally:
   ```bash
   curl http://localhost:8080
   ```

## Backup and Recovery

### Backup

1. Database:
   ```bash
   cp <INSTALL_DIR>/backend/test.db ~/backup/
   ```

2. Configuration:
   ```bash
   cp ~/.unity/unity.yaml ~/backup/
   ```

3. Module Registry (if using local):
   ```bash
   cp <WORKDIR>/module-registry.json ~/backup/
   ```

4. Terraform state (automatically in S3)

### Recovery

1. Restore configuration file
2. Restore database
3. Restart service

## Updating the Console

```bash
cd <INSTALL_DIR>
git pull origin main
cd backend
go build -o management-console cmd/web/main.go
cd ../ui
npm install
npm run build
sudo systemctl restart unity-management-console
```

## Alternative Deployment Options

For containerized deployments that don't require constant operation, consider the **ECS Fargate deployment option** which provides:
- On-demand scaling (start/stop as needed)
- Persistent storage with EFS
- No server management
- Cost-effective pay-per-use model

See the [ECS Fargate Deployment Guide](ecs-fargate-deployment-guide.md) for detailed instructions.

## Additional Resources

- [Unity SDS Documentation](https://unity-sds.gitbook.io/docs/)
- [Unity Marketplace](https://github.com/unity-sds/unity-marketplace)
- [Management Console Repository](https://github.com/unity-sds/unity-management-console)
- [ECS Fargate Deployment Guide](ecs-fargate-deployment-guide.md)
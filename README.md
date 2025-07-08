# Unity Management Console

The Unity Management Console is a comprehensive web-based management portal for deploying and managing Unity Science Data System (SDS) applications. It provides a unified interface for infrastructure provisioning, application lifecycle management, and system configuration.

## Architecture Overview

The Unity Management Console follows a modular architecture designed for managing cloud-native applications:

### Core Components

- **Backend API Server**: Go-based REST API using Gin framework
- **Frontend UI**: Svelte-based single-page application
- **Database**: SQLite with GORM ORM for application state and audit logging
- **Infrastructure as Code**: Terraform integration for resource provisioning
- **Real-time Updates**: WebSocket connections for live status updates
- **Marketplace Integration**: Package management system for Unity applications

### Key Technologies

- **Backend**: Go, Gin Web Framework, GORM, Gorilla WebSocket
- **Frontend**: Svelte, TypeScript, Tailwind CSS
- **Infrastructure**: Terraform, AWS SDK
- **Database**: SQLite (local), S3 (Terraform state), DynamoDB (state locking)

## Software Capabilities

### Application Management
- **Install Applications**: Deploy Unity applications from the marketplace with version control
- **Update Applications**: Seamlessly upgrade applications to newer versions
- **Uninstall Applications**: Clean removal of applications and their resources
- **Application Health Monitoring**: Real-time status checks and health indicators

### Infrastructure Management
- **Automated Provisioning**: Terraform-based infrastructure deployment
- **State Management**: Centralized Terraform state with S3 backend and DynamoDB locking
- **Resource Tracking**: Monitor and manage AWS resources created by applications

### Configuration Management
- **SSM Parameter Store Integration**: Secure storage and retrieval of configuration parameters
- **Environment Customization**: Project and venue-specific configurations
- **Bootstrap Process**: Automated initial setup of core Unity components

### Marketplace Features
- **Application Discovery**: Browse available Unity applications and their versions
- **Metadata Management**: Access application descriptions, dependencies, and requirements
- **Version Control**: Install specific versions of applications
- **Custom Marketplace Support**: Configure alternative marketplace repositories

### Operational Features
- **Audit Logging**: Track all installation, update, and configuration changes
- **Installation Logs**: Detailed logs for troubleshooting deployments
- **Batch Operations**: Install multiple applications in a single Terraform run
- **Pre/Post Install Scripts**: Execute custom scripts during deployment

### Security Features
- **IAM Integration**: AWS IAM-based access control
- **Secure Parameter Storage**: Sensitive data stored in AWS SSM Parameter Store
- **HTTPS Support**: TLS encryption for web interface (when configured with reverse proxy)

## Deployment Options

- **Local Development**: Run directly on developer machines
- **EC2 Deployment**: Production deployment on AWS EC2 instances (see [EC2 Deployment Guide](documentation/ec2-deployment-guide.md))
- **ECS Fargate Deployment**: Containerized deployment on AWS ECS Fargate with EFS storage (see [ECS Fargate Deployment Guide](documentation/ecs-fargate-deployment-guide.md))
- **Container Deployment**: Docker support for containerized environments
- **Systemd Service**: Run as a system service on Linux

## Module Registry

The Unity Management Console supports a Module Registry that allows you to provide reusable Terraform modules for applications to use. This promotes consistency and best practices across deployments. See the [Module Registry Guide](documentation/module-registry-guide.md) for detailed information on:

- Setting up a module registry
- Defining reusable Terraform modules
- Referencing modules in application installations
- Version management and caching

## Configuration

### Environment Variables

Local testing and development requires the following environment variables:

```bash
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_SESSION_TOKEN
GITHUB_TOKEN
```

### Configuration File

The Unity Management Console looks for configuration in the following order:
1. Path specified by `--config` CLI flag
2. Path specified by `UNITY_CONFIG_PATH` environment variable  
3. Default location: `~/.unity/unity.yaml`

Create a configuration file (e.g., `~/.unity/unity.yaml`):

```yaml
# AWS Configuration
awsregion: "us-west-2"
bucketname: "unity-<project>-<venue>-terraform-state"
project: "myproject"
venue: "dev"
installprefix: "unity"

# Console Configuration
consolehost: "http://localhost:8080"
basepath: "http://localhost"
workdir: "./workdir"

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
```

### Module Registry (Optional)

The Unity Management Console supports a Module Registry for providing reusable Terraform modules. To enable this feature:

1. **Option 1 - Local Registry**: Place a `module-registry.json` file in your workdir:
   ```
   {workdir}/module-registry.json
   ```

2. **Option 2 - Remote Registry**: The registry will be automatically downloaded from your marketplace repository:
   ```
   https://raw.githubusercontent.com/{owner}/{repo}/main/module-registry.json
   ```

See the [Module Registry Guide](documentation/module-registry-guide.md) for detailed setup instructions and examples.

## Running
Grab a management console zip file from the [releases page](https://github.com/unity-sds/unity-management-console/releases)
Unzip to your target destination.

```shell

unzip managementconsole.zip
cd management-console
./main
```

## Development

```shell
npm install
```
Make sure dependencies are installed.

```shell
npm run dev
```
Run the development environment. This is a Svelte only environment it provides a fake backend and http responses.

```shell
npm run build-all
```
Build the frontend and the backend. This also runs various linters and other code quality checks.

```shell
npm run serve
```
Launch the frontend and the backend locally for testing.



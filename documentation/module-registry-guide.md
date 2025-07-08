# Unity Management Console Module Registry Guide

This guide explains how to use the Module Registry feature in the Unity Management Console to provide reusable Terraform modules that can be used directly in Terraform configurations via the Terraform Registry Protocol.

## Overview

The Unity Management Console acts as a **Terraform Module Registry**, allowing you to:

- Define and catalog reusable Terraform modules
- Reference modules cleanly in Terraform code using standard registry syntax
- Manage module versions and dependencies
- Cache modules locally for performance

## How It Works

1. **Module Registry**: The Management Console loads a module registry configuration
2. **Terraform Registry API**: Exposes standard Terraform Registry Protocol endpoints
3. **Standard Terraform Workflow**: Use modules with normal `terraform` commands
4. **Automatic Caching**: Modules are downloaded and cached automatically

## Setting Up the Module Registry

### Option 1: Local Registry File

Place a `module-registry.json` file in your workdir:

```bash
~/.unity/workdir/module-registry.json
```

### Option 2: Marketplace Registry

The module registry will be automatically downloaded from:
```
https://raw.githubusercontent.com/{owner}/{repo}/main/module-registry.json
```

Based on your marketplace configuration in your Unity config file (default: `~/.unity/unity.yaml`).

## Module Development Requirements

### Module Structure Requirements

Modules in the registry must follow these standards:

#### 1. Git Repository Structure
```
terraform-module-vpc/
├── main.tf           # Primary module logic
├── variables.tf      # Input variable definitions
├── outputs.tf        # Output value definitions  
├── README.md         # Module documentation
├── versions.tf       # Provider requirements
└── examples/         # Usage examples (recommended)
    └── basic/
        ├── main.tf
        └── variables.tf
```

#### 2. Git Tagging for Versions
- **Required**: Use semantic versioning tags
- **Format**: `v1.0.0` or `1.0.0` (both supported)
- **Examples**: `v1.0.0`, `v1.1.0`, `v2.0.0`
- **Special**: `latest` version points to default branch

#### 3. Provider Requirements
Include `versions.tf` with provider constraints:
```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
  required_version = ">= 1.0"
}
```

#### 4. Variable and Output Documentation
**variables.tf:**
```hcl
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  validation {
    condition     = can(cidrhost(var.vpc_cidr, 0))
    error_message = "VPC CIDR must be a valid IPv4 CIDR block."
  }
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b"]
}
```

**outputs.tf:**
```hcl
output "vpc_id" {
  description = "ID of the created VPC"
  value       = aws_vpc.main.id
}

output "private_subnet_ids" {
  description = "IDs of private subnets"
  value       = aws_subnet.private[*].id
}
```

### Versioning Standards

#### Semantic Versioning Rules
- **MAJOR**: Breaking changes (e.g., removed variables, changed outputs)
- **MINOR**: New features, new optional variables
- **PATCH**: Bug fixes, documentation updates

#### Git Workflow
```bash
# Release a new version
git tag v1.2.0
git push origin v1.2.0

# Update latest pointer (optional)
git tag -f latest
git push origin latest --force
```

## Module Registry Format

### Complete module-registry.json Schema

Based on the Unity Management Console code, here's the complete schema with required and optional fields:

```json
{
  // REQUIRED: Registry metadata
  "version": "1.0",                           // REQUIRED: Registry format version
  "metadata": {                              // REQUIRED: Registry metadata
    "name": "Unity Terraform Modules",       // REQUIRED: Registry name
    "description": "Module registry description", // REQUIRED: Registry description
    "last_updated": "2025-06-30T10:00:00Z"  // REQUIRED: ISO timestamp
  },
  
  // REQUIRED: Modules collection
  "modules": {
    "infrastructure-vpc": {                  // REQUIRED: Module key (namespace-name format)
      
      // REQUIRED FIELDS
      "description": "Standard Unity VPC module",  // REQUIRED: Module description
      "provider": "unity",                         // REQUIRED: Provider name (e.g., "unity", "aws", "google")
      "source": "github.com/unity-sds/unity-terraform-modules//modules/vpc", // REQUIRED: Git source
      "versions": {                          // REQUIRED: Version definitions
        "1.0.0": {                          // Version identifier
          "ref": "v1.0.0"                   // REQUIRED: Git reference (tag/branch)
        },
        "1.1.0": {
          "ref": "v1.1.0",
          // OPTIONAL VERSION FIELDS:
          "terraform_version": ">= 1.0",    // OPTIONAL: Terraform version constraint
          "min_mc_version": "1.2.0"        // OPTIONAL: Minimum Management Console version
        },
        "latest": {
          "ref": "main"                     // Points to default branch
        }
      },
      
      // OPTIONAL FIELDS
      "documentation": "https://docs.unity.nasa.gov/modules/vpc", // OPTIONAL: Documentation URL
      "inputs": {                          // OPTIONAL: Input variable definitions
        "vpc_cidr": {
          "description": "CIDR block for VPC",     // REQUIRED if inputs defined
          "type": "string",                        // REQUIRED if inputs defined  
          "required": true,                        // REQUIRED if inputs defined
          "default": "10.0.0.0/16",               // OPTIONAL: Default value
          "sensitive": false                       // OPTIONAL: Mark as sensitive
        },
        "availability_zones": {
          "description": "List of availability zones",
          "type": "list(string)",
          "required": false,
          "default": ["us-west-2a", "us-west-2b"]
        }
      },
      "outputs": {                         // OPTIONAL: Output definitions  
        "vpc_id": {
          "description": "ID of the created VPC"     // REQUIRED if outputs defined
        },
        "private_subnet_ids": {
          "description": "IDs of private subnets"
        }
      }
    }
  }
}
```

### Required vs Optional Fields Summary

#### **Registry Level (REQUIRED)**
- `version` - Registry format version
- `metadata` - Registry metadata object
  - `name` - Registry name
  - `description` - Registry description  
  - `last_updated` - ISO timestamp
- `modules` - Modules collection

#### **Module Level**
**REQUIRED:**
- Module key (e.g., `"infrastructure-vpc"`)
- `description` - Module description
- `provider` - Provider name (e.g., `"unity"`, `"aws"`, `"google"`, `"azure"`)
- `source` - Git repository source
- `versions` - Version definitions object
  - Each version must have `ref` (Git reference)

**OPTIONAL:**
- `documentation` - Documentation URL
- `inputs` - Input variable definitions
- `outputs` - Output definitions

#### **Version Level**
**REQUIRED:**
- `ref` - Git tag or branch reference

**OPTIONAL:**
- `terraform_version` - Terraform version constraints
- `min_mc_version` - Minimum Management Console version

#### **Input Level (if inputs defined)**
**REQUIRED:**
- `description` - Input description
- `type` - Terraform type (string, list(string), etc.)
- `required` - Boolean indicating if required

**OPTIONAL:**
- `default` - Default value
- `sensitive` - Mark as sensitive (boolean)

#### **Output Level (if outputs defined)**
**REQUIRED:**
- `description` - Output description

### Provider Requirements

Provider requirements are **only** defined in the module's `versions.tf` file, not in the registry metadata.

#### **Module `versions.tf` (Required)**
Every module must include a `versions.tf` file with provider requirements:

```hcl
# In the actual module repository at the git ref
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes" 
      version = ">= 2.0"
    }
  }
}
```

**Purpose**:
- **Enforced by Terraform** during `terraform init`
- **Required for module to work** - Terraform reads this from the downloaded module
- **Source of truth** for all provider requirements

**Why no registry metadata for providers?**
- Avoids duplication and potential inconsistencies
- Terraform already provides this information from the module itself
- Keeps registry metadata minimal and focused

### Minimal Module Registry Example

```json
{
  "version": "1.0",
  "metadata": {
    "name": "Minimal Registry",
    "description": "Basic example",
    "last_updated": "2025-01-01T00:00:00Z"
  },
  "modules": {
    "infrastructure-vpc": {
      "description": "VPC module",
      "provider": "unity",
      "source": "github.com/org/terraform-module-vpc",
      "versions": {
        "1.0.0": {
          "ref": "v1.0.0"
        }
      }
    }
  }
}
```

This minimal example shows only the **required** fields - the module will work as long as the Git repository contains proper Terraform code with `versions.tf`.

## Using Modules in Terraform

Once the Management Console is running with a module registry, you can reference modules directly in your Terraform configurations:

### Basic Module Reference

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Basic module reference
module "vpc" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "~> 1.0"
  
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
  name_prefix        = "my-project"
}
```

### Version Constraints

```hcl
# Exact version
module "vpc" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "1.2.3"
  # ...
}

# Pessimistic constraint (recommended)
module "vpc" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "~> 1.2"    # >= 1.2.0, < 2.0.0
  # ...
}

# Range constraint
module "vpc" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = ">= 1.0, < 2.0"
  # ...
}

# Latest version (not recommended for production)
module "vpc" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "latest"
  # ...
}
```

### Module Composition and Dependencies

```hcl
# VPC module
module "vpc" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "~> 1.0"
  
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
  name_prefix        = var.project_name
}

# Database module that depends on VPC
module "database" {
  source  = "my-unity-console.com/infrastructure/rds/unity"
  version = "~> 2.1"
  
  # Reference outputs from VPC module
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnet_ids
  security_group_ids = [module.vpc.database_security_group_id]
  
  db_name     = "${var.project_name}_db"
  db_username = var.db_username
  db_password = var.db_password
}

# Application module that depends on both VPC and Database
module "application" {
  source  = "my-unity-console.com/infrastructure/ecs-app/unity"
  version = "~> 3.0"
  
  vpc_id         = module.vpc.vpc_id
  subnet_ids     = module.vpc.private_subnet_ids
  database_url   = module.database.connection_string
  
  app_name       = var.project_name
  app_image      = var.app_image
  desired_count  = var.app_desired_count
}
```

### Environment-Specific Configurations

```hcl
# Development environment
module "vpc_dev" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "~> 1.0"
  
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-west-2a", "us-west-2b"]
  environment        = "dev"
  
  # Smaller NAT gateway for cost savings
  single_nat_gateway = true
}

# Production environment  
module "vpc_prod" {
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "~> 1.0"
  
  vpc_cidr           = "10.1.0.0/16"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
  environment        = "prod"
  
  # High availability setup
  single_nat_gateway = false
  enable_flow_logs   = true
}
```

### Advanced Usage with Count and For_Each

```hcl
# Multiple environments using count
module "vpc" {
  count   = length(var.environments)
  source  = "my-unity-console.com/infrastructure/vpc/unity"
  version = "~> 1.0"
  
  vpc_cidr    = var.vpc_cidrs[count.index]
  environment = var.environments[count.index]
  name_prefix = "${var.project_name}-${var.environments[count.index]}"
}

# Multiple regions using for_each
module "vpc_multi_region" {
  for_each = var.regions
  source   = "my-unity-console.com/infrastructure/vpc/unity"
  version  = "~> 1.0"
  
  providers = {
    aws = aws.region[each.key]
  }
  
  vpc_cidr    = each.value.cidr
  region      = each.key
  name_prefix = "${var.project_name}-${each.key}"
}

### Testing and Validation

```bash
# Initialize and download modules
terraform init

# Validate syntax and configuration
terraform validate

# Format code consistently
terraform fmt

# Plan your infrastructure (dry run)
terraform plan

# Apply changes
terraform apply

# Check what modules are installed
terraform providers

# Update modules to latest versions
terraform init -upgrade
```

### Debugging Module Issues

```bash
# Enable detailed logging
export TF_LOG=DEBUG
terraform init

# Check module cache
ls -la .terraform/modules/

# Manually verify module source
curl -I https://my-unity-console.com/v1/modules/unity/vpc/aws/versions

# Test specific module version
terraform init -upgrade=true
```

### Local Development and Testing

For local development, you can override module sources:

```hcl
# Override for local testing
module "vpc" {
  source = "../local-modules/vpc"  # Local override
  # source  = "my-unity-console.com/unity/vpc/aws"  # Production
  # version = "~> 1.0"
  
  vpc_cidr = "10.0.0.0/16"
}
```

## Registry API Endpoints

The Management Console exposes these **standard Terraform Registry Protocol** endpoints:

- `GET /.well-known/terraform.json` - Service discovery
- `GET /v1/modules/{namespace}/{name}/{provider}/versions` - List module versions
- `GET /v1/modules/{namespace}/{name}/{provider}/{version}/download` - Get download URL
- `GET /v1/modules/{namespace}/{name}/{provider}/{version}/archive` - Download module archive

**Important**: These paths are **required by the Terraform Registry Protocol**. Terraform automatically appends these paths when resolving module sources.

### How Terraform Resolves Modules

When you write:
```hcl
module "vpc" {
  source = "my-console.com/infrastructure/vpc/unity"
}
```

Terraform automatically:
1. Calls `https://my-console.com/.well-known/terraform.json` (service discovery)
2. Calls `https://my-console.com/v1/modules/infrastructure/vpc/unity/versions` (get versions)
3. Calls `https://my-console.com/v1/modules/infrastructure/vpc/unity/{version}/download` (get download URL)

### Example API Calls

```bash
# Service discovery (Terraform does this automatically)
curl https://my-console.com/.well-known/terraform.json

# List versions for infrastructure/vpc/unity module
curl https://my-console.com/v1/modules/infrastructure/vpc/unity/versions

# Get download URL for specific version
curl https://my-console.com/v1/modules/infrastructure/vpc/unity/1.0.0/download
```

## Module Naming and Path Construction

### Understanding the Provider Field

The **provider field is REQUIRED** in the module registry and determines the last part of the Terraform module source URL:

```json
{
  "modules": {
    "infrastructure-vpc": {
      "provider": "unity",    // This becomes the provider in the URL
      // ...
    }
  }
}
```

This results in the Terraform source: `my-console.com/infrastructure/vpc/unity`

### Path Construction Flow

```
Git Repository: terraform-module-vpc
       ↓
Registry Entry: "infrastructure-vpc" with provider: "unity"
       ↓
Terraform Reference: my-console.com/infrastructure/vpc/unity
```

### Naming Convention Rules

1. **Git Repository Name**: Can be anything (e.g., `terraform-module-vpc`, `vpc-module`, `infrastructure-vpc`)
2. **Registry Entry Key**: `{namespace}-{name}` format (e.g., `"infrastructure-vpc"`)
   - **Namespace**: Extracted as first part before `-` (e.g., `infrastructure`)
   - **Name**: Extracted as second part after `-` (e.g., `vpc`)
3. **Provider**: Explicitly defined in the `provider` field (e.g., `"unity"`, `"aws"`, `"google"`)
4. **Terraform Module Reference**: `{mc-domain}/{namespace}/{name}/{provider}`

### Provider Field Purpose

The provider field serves several purposes:

1. **Terraform Registry Protocol Compliance**: The Terraform Registry Protocol requires URLs in the format `/v1/modules/{namespace}/{name}/{provider}/versions`
2. **Module Organization**: Allows the same module name to exist for different providers
3. **Clear Module Identity**: Makes it explicit which provider ecosystem the module belongs to

### Common Provider Values

- `"unity"` - Unity-specific modules
- `"aws"` - AWS-specific modules
- `"google"` - Google Cloud modules
- `"azure"` - Azure modules
- `"kubernetes"` - Kubernetes modules
- Custom provider names for your organization

### Example Mapping

| Component | Value | 
|-----------|-------|
| Git Repository | `terraform-module-vpc` |
| Registry Key | `"infrastructure-vpc"` |
| Namespace | `infrastructure` (first part before `-`) |
| Name | `vpc` (second part after `-`) |
| Provider | `unity` (from provider field) |
| **Terraform Reference** | `my-console.com/infrastructure/vpc/unity` |

### Multiple Provider Examples

```json
{
  "modules": {
    "infrastructure-vpc": {
      "description": "Unity VPC module",
      "provider": "unity",
      "source": "github.com/unity-sds/terraform-modules//modules/vpc"
      // Terraform: my-console.com/infrastructure/vpc/unity
    },
    "infrastructure-vpc": {
      "description": "AWS VPC module",
      "provider": "aws",
      "source": "github.com/aws-modules/terraform-modules//modules/vpc"
      // Terraform: my-console.com/infrastructure/vpc/aws
    },
    "infrastructure-storage": {
      "description": "Google Cloud Storage module",
      "provider": "google",
      "source": "github.com/gcp-modules/terraform-modules//modules/storage"
      // Terraform: my-console.com/infrastructure/storage/google
    },
    "company-database": {
      "description": "Company RDS module",
      "provider": "unity",
      "source": "github.com/company/terraform-modules//modules/rds"
      // Terraform: my-console.com/company/database/unity
    }
  }
}
```

**Important**: The provider field in the module registry determines the URL path. This is separate from the Terraform providers defined in the module's `versions.tf` file.

## Benefits

### For Module Developers
- Version control and release management
- Centralized documentation and discovery
- Consistent module interface standards

### For Infrastructure Users  
- Clean, readable Terraform code
- Standard Terraform tooling works unchanged
- Version constraints and dependency management
- Local IDE support and validation

### For Organizations
- Enforce infrastructure standards and best practices
- Promote code reuse across teams
- Centralized governance and compliance
- Audit trail of module usage

## Migration from Application moduleReferences

If you were previously using `moduleReferences` in application installation parameters, you should migrate to using modules directly in your Terraform code:

**Before (moduleReferences approach - no longer supported):**
```json
{
  "moduleReferences": [
    {
      "name": "unity-vpc",
      "alias": "main_vpc",
      "config": {
        "vpc_cidr": "10.0.0.0/16"
      }
    }
  ]
}
```

**After (Terraform registry approach):**
```hcl
module "main_vpc" {
  source  = "my-console.com/unity/vpc/aws"
  version = "~> 1.0"
  
  vpc_cidr = "10.0.0.0/16"
}
```

This new approach provides better integration with standard Terraform workflows and tooling.
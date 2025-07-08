# Unity Terraform Module Development Guide

This guide covers best practices for developing Terraform modules that work with the Unity Management Console's module registry.

## Module Development Standards

### Repository Structure

```
terraform-module-{name}/
├── main.tf                 # Primary module logic
├── variables.tf            # Input variable definitions
├── outputs.tf              # Output value definitions
├── versions.tf             # Provider and Terraform version constraints
├── README.md               # Module documentation
├── CHANGELOG.md            # Version history
├── LICENSE                 # Module license
├── .gitignore              # Git ignore patterns
├── examples/               # Usage examples
│   ├── basic/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── terraform.tfvars.example
│   └── advanced/
│       ├── main.tf
│       ├── variables.tf
│       └── terraform.tfvars.example
└── tests/                  # Automated tests (optional)
    ├── basic_test.go
    └── advanced_test.go
```

### File Conventions

#### main.tf
- Contains primary resource definitions
- Groups related resources logically
- Uses data sources when appropriate
- Includes local values for computed values

```hcl
# main.tf
locals {
  common_tags = {
    Project     = var.project
    Environment = var.environment
    ManagedBy   = "terraform"
    Module      = "unity-vpc"
  }
}

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = var.enable_dns_hostnames
  enable_dns_support   = var.enable_dns_support
  
  tags = merge(local.common_tags, {
    Name = "${var.name_prefix}-vpc"
  })
}

resource "aws_subnet" "private" {
  count = length(var.availability_zones)
  
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone = var.availability_zones[count.index]
  
  tags = merge(local.common_tags, {
    Name = "${var.name_prefix}-private-${var.availability_zones[count.index]}"
    Type = "private"
  })
}
```

#### variables.tf
- All inputs with descriptions and types
- Include validation blocks where appropriate
- Set sensible defaults for optional variables

```hcl
# variables.tf
variable "project" {
  description = "Unity project name"
  type        = string
}

variable "environment" {
  description = "Deployment environment (dev, test, prod)"
  type        = string
  
  validation {
    condition     = contains(["dev", "test", "prod"], var.environment)
    error_message = "Environment must be one of: dev, test, prod."
  }
}

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
  
  validation {
    condition     = length(var.availability_zones) >= 2
    error_message = "At least 2 availability zones must be specified."
  }
}

variable "enable_dns_hostnames" {
  description = "Enable DNS hostnames in VPC"
  type        = bool
  default     = true
}

variable "enable_dns_support" {
  description = "Enable DNS support in VPC"
  type        = bool
  default     = true
}

variable "name_prefix" {
  description = "Name prefix for all resources"
  type        = string
}

variable "tags" {
  description = "Additional tags to apply to resources"
  type        = map(string)
  default     = {}
}
```

#### outputs.tf
- All important resource attributes
- Include descriptions for all outputs
- Group related outputs logically

```hcl
# outputs.tf
output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = aws_subnet.private[*].id
}

output "private_subnet_cidrs" {
  description = "CIDR blocks of the private subnets"
  value       = aws_subnet.private[*].cidr_block
}

output "availability_zones" {
  description = "Availability zones used"
  value       = var.availability_zones
}

output "tags" {
  description = "Tags applied to resources"
  value       = local.common_tags
}
```

#### versions.tf
- Specify minimum provider versions
- Set minimum Terraform version
- Be conservative with version constraints

```hcl
# versions.tf
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
}
```

### Documentation Standards

#### README.md Template

```markdown
# Unity VPC Terraform Module

This module creates a VPC with private subnets suitable for Unity applications.

## Usage

```hcl
module "vpc" {
  source  = "my-console.com/unity/vpc/aws"
  version = "~> 1.0"
  
  project            = "my-project"
  environment        = "prod"
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
  name_prefix        = "unity-prod"
}
```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| aws | >= 4.0 |

## Providers

| Name | Version |
|------|---------|
| aws | >= 4.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| project | Unity project name | `string` | n/a | yes |
| environment | Deployment environment | `string` | n/a | yes |
| vpc_cidr | CIDR block for VPC | `string` | n/a | yes |
| availability_zones | List of availability zones | `list(string)` | `["us-west-2a", "us-west-2b"]` | no |
| name_prefix | Name prefix for all resources | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| vpc_id | ID of the VPC |
| private_subnet_ids | IDs of the private subnets |
| availability_zones | Availability zones used |

## Examples

See the [examples](./examples/) directory for usage examples.
```

### Version Management

#### Semantic Versioning
- **MAJOR.MINOR.PATCH** format
- **MAJOR**: Breaking changes (removed variables, changed outputs, changed behavior)
- **MINOR**: New features, new optional variables, new outputs
- **PATCH**: Bug fixes, documentation updates, internal improvements

#### Git Tagging
```bash
# Create and push version tags
git tag v1.0.0
git push origin v1.0.0

# Update latest tag (optional)
git tag -f latest
git push origin latest --force
```

#### CHANGELOG.md Format
```markdown
# Changelog

## [1.2.0] - 2025-01-15

### Added
- New `enable_flow_logs` variable for VPC flow logs
- NAT Gateway configuration options

### Changed
- Default subnet size changed from /24 to /20

### Fixed
- Fixed route table associations for private subnets

## [1.1.0] - 2025-01-01

### Added
- Support for multiple availability zones
- Custom tags input variable
```

### Testing Best Practices

#### Example Testing Structure
```hcl
# examples/basic/main.tf
module "vpc" {
  source = "../../"
  
  project            = "test-project"
  environment        = "dev"
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-west-2a", "us-west-2b"]
  name_prefix        = "test"
}

output "vpc_id" {
  value = module.vpc.vpc_id
}
```

#### Manual Testing Checklist
- [ ] `terraform init` succeeds
- [ ] `terraform validate` passes
- [ ] `terraform plan` shows expected resources
- [ ] `terraform apply` creates resources successfully
- [ ] All outputs are populated correctly
- [ ] `terraform destroy` cleans up all resources
- [ ] Examples work correctly

### Module Registry Configuration

#### Registry Entry Example
```json
{
  "unity-vpc": {
    "description": "Standard Unity VPC module with private subnets",
    "source": "github.com/unity-sds/terraform-module-vpc",
    "provider": "hashicorp/aws",
    "documentation": "https://github.com/unity-sds/terraform-module-vpc",
    "versions": {
      "1.0.0": {
        "ref": "v1.0.0",
        "min_provider_version": ">= 4.0"
      },
      "1.1.0": {
        "ref": "v1.1.0",
        "min_provider_version": ">= 4.0"
      },
      "latest": {
        "ref": "main"
      }
    },
    "inputs": {
      "project": {
        "description": "Unity project name",
        "type": "string",
        "required": true
      },
      "environment": {
        "description": "Deployment environment",
        "type": "string",
        "required": true
      },
      "vpc_cidr": {
        "description": "CIDR block for VPC",
        "type": "string",
        "required": true
      },
      "availability_zones": {
        "description": "List of availability zones",
        "type": "list(string)",
        "required": false,
        "default": ["us-west-2a", "us-west-2b"]
      }
    },
    "outputs": {
      "vpc_id": {
        "description": "ID of the VPC"
      },
      "private_subnet_ids": {
        "description": "IDs of the private subnets"
      }
    }
  }
}
```

## Common Patterns

### Data Sources
```hcl
# Use data sources for external dependencies
data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}
```

### Conditional Resources
```hcl
# Conditional resource creation
resource "aws_nat_gateway" "main" {
  count = var.enable_nat_gateway ? length(var.availability_zones) : 0
  
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  
  tags = merge(local.common_tags, {
    Name = "${var.name_prefix}-nat-${var.availability_zones[count.index]}"
  })
}
```

### Dynamic Blocks
```hcl
# Dynamic configuration blocks
resource "aws_security_group" "main" {
  name_prefix = "${var.name_prefix}-"
  vpc_id      = aws_vpc.main.id
  
  dynamic "ingress" {
    for_each = var.security_group_rules
    content {
      from_port   = ingress.value.from_port
      to_port     = ingress.value.to_port
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.cidr_blocks
    }
  }
}
```

## Quality Checklist

### Before Release
- [ ] All variables have descriptions and appropriate types
- [ ] All outputs have descriptions
- [ ] Provider versions are specified
- [ ] Examples are working and up-to-date
- [ ] README.md is complete and accurate
- [ ] CHANGELOG.md is updated
- [ ] Git tag follows semantic versioning
- [ ] Manual testing completed
- [ ] Module follows Unity naming conventions

### Security Considerations
- [ ] No hardcoded secrets or sensitive data
- [ ] Appropriate use of sensitive = true for outputs
- [ ] Security groups follow least privilege principle
- [ ] IAM policies are minimal and specific
- [ ] Resources are properly tagged for cost tracking
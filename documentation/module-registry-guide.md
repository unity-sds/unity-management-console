# Unity Management Console Module Registry Guide

This guide explains how to use the Module Registry feature in the Unity Management Console to provide reusable Terraform modules that applications can reference during installation.

## Overview

The Module Registry allows the Unity Management Console to "own" and provide standardized Terraform modules that applications in the catalog can use. This promotes code reuse, consistency, and best practices across Unity deployments.

## How It Works

1. **Module Registry File**: A JSON catalog defines available modules, their versions, inputs, and outputs
2. **Module References**: Applications can reference modules from the registry during installation
3. **Automatic Caching**: Modules are automatically downloaded and cached locally
4. **Version Control**: Specific versions or "latest" can be referenced

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

Based on your marketplace configuration in `~/.unity/unity.yaml`.

## Module Registry Format

```json
{
  "version": "1.0",
  "metadata": {
    "name": "Unity Terraform Modules",
    "description": "Module registry description",
    "last_updated": "2025-06-30T10:00:00Z"
  },
  "modules": {
    "module-name": {
      "description": "Module description",
      "source": "github.com/org/repo//path/to/module",
      "documentation": "https://docs.example.com",
      "versions": {
        "1.0.0": {
          "ref": "v1.0.0",
          "terraform_version": ">= 1.0",
          "providers": {
            "aws": ">= 4.0"
          }
        },
        "latest": {
          "ref": "main"
        }
      },
      "inputs": {
        "variable_name": {
          "description": "Variable description",
          "type": "string",
          "required": true,
          "default": "default_value"
        }
      },
      "outputs": {
        "output_name": "Output description"
      }
    }
  }
}
```

## Using Modules in Applications

Applications can reference modules in their installation parameters:

```json
{
  "name": "my-application",
  "version": "1.0.0",
  "deploymentName": "my-app-prod",
  "moduleReferences": [
    {
      "name": "unity-vpc",
      "version": "1.1.0",
      "alias": "main_vpc",
      "config": {
        "vpc_cidr": "10.0.0.0/16",
        "name_prefix": "unity-prod",
        "availability_zones": ["us-west-2a", "us-west-2b"]
      }
    }
  ]
}
```

### Module Reference Fields

- **name**: Module name from the registry (required)
- **version**: Module version or "latest" (optional, defaults to "latest")
- **alias**: Custom name for the module instance (optional)
- **config**: Module input variables (required based on module definition)
- **depends_on**: List of other module aliases this depends on (optional)

## Module Resolution Process

1. **Registry Loading**: The console loads the module registry on startup
2. **Module Resolution**: When an application references a module, it's resolved from the registry
3. **Validation**: Required inputs are validated against the module definition
4. **Caching**: The module is downloaded to `{workdir}/module_cache/`
5. **Terraform Generation**: Module blocks are added to the generated Terraform configuration

## Generated Terraform

Module references are converted to standard Terraform module blocks:

```hcl
module "main_vpc" {
  source = "/home/user/.unity/workdir/module_cache/unity-vpc-1.1.0"
  
  vpc_cidr           = "10.0.0.0/16"
  name_prefix        = "unity-prod"
  availability_zones = ["us-west-2a", "us-west-2b"]
  
  # Auto-injected from app config
  project      = "myproject"
  venue        = "prod"
  installprefix = "unity"
}
```

## Cross-Module References

Modules can reference outputs from other modules:

```json
{
  "name": "unity-rds",
  "version": "latest",
  "alias": "database",
  "config": {
    "vpc_id": "${module.main_vpc.vpc_id}",
    "subnet_ids": "${module.main_vpc.private_subnet_ids}"
  },
  "depends_on": ["main_vpc"]
}
```

## Backward Compatibility

This feature is fully backward compatible:

- If no module registry exists, the console operates normally
- Applications without `moduleReferences` work as before
- Module references are optional and don't affect existing functionality

## Best Practices

1. **Version Pinning**: Use specific versions in production
2. **Module Naming**: Use descriptive, consistent module names
3. **Documentation**: Keep module documentation up to date
4. **Testing**: Test modules thoroughly before adding to registry
5. **Inputs**: Define sensible defaults for optional inputs

## Troubleshooting

### Module Not Found

Check that:
- Module registry file exists and is valid JSON
- Module name matches exactly
- Registry URL is accessible (if using remote registry)

### Version Not Found

Ensure:
- Requested version exists in the module's versions map
- Git ref (tag/branch) exists in the source repository

### Module Download Fails

Verify:
- Network connectivity to the module source
- Git is installed and available in PATH
- Proper authentication for private repositories

## Example Modules

See the `examples/module-registry.json` file for a complete example registry with common Unity modules:

- **unity-vpc**: Standard VPC with public/private subnets
- **unity-eks**: EKS cluster configuration
- **unity-rds**: RDS database instances
- **unity-s3**: S3 buckets with standard settings

## Security Considerations

1. **Source Validation**: Only use modules from trusted sources
2. **Version Control**: Pin versions to avoid unexpected changes
3. **Access Control**: Ensure module repositories have proper access controls
4. **Sensitive Data**: Never hardcode secrets in module configurations
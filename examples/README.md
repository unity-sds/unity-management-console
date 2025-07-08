# Unity Management Console Examples

This directory contains example configurations and files for the Unity Management Console.

## Files

### module-registry.json
A complete example of a module registry that provides reusable Terraform modules for Unity applications. This includes:
- unity-vpc: VPC with public/private subnets
- unity-eks: EKS cluster configuration
- unity-rds: RDS PostgreSQL database
- unity-s3: S3 bucket with standard settings

### install-with-modules.json
Example of an application installation request that references modules from the registry. Shows how to:
- Reference specific module versions
- Configure module inputs
- Set up dependencies between modules
- Use module outputs in other modules

### unity-config-with-modules.yaml
Complete Unity Management Console configuration file example that includes:
- Basic AWS and project settings
- Marketplace configuration
- Bootstrap applications
- SSM parameters
- Module registry configuration notes

## Usage

### Setting Up Module Registry

1. Copy `module-registry.json` to your workdir:
   ```bash
   cp module-registry.json ~/.unity/workdir/
   ```

2. Or commit it to your marketplace repository at the root level.

### Installing Applications with Modules

Use the `install-with-modules.json` as a template for your application installation requests:

```bash
curl -X POST http://localhost:8080/api/applications/install \
  -H "Content-Type: application/json" \
  -d @install-with-modules.json
```

### Configuration

Use `unity-config-with-modules.yaml` as a starting point for your Unity configuration:

```bash
# Copy to default location
cp unity-config-with-modules.yaml ~/.unity/unity.yaml

# Or specify custom location with --config flag:
# ./unity webapp --config ./unity-config-with-modules.yaml

# Or set environment variable:
# export UNITY_CONFIG_PATH=./unity-config-with-modules.yaml
# ./unity webapp

# Edit the file to match your environment
```

## Module Development

When creating new modules for the registry:

1. Follow Terraform module best practices
2. Use semantic versioning for releases
3. Document all inputs and outputs
4. Test modules thoroughly before adding to registry
5. Use git tags for version references

## Additional Resources

- [Module Registry Guide](../documentation/module-registry-guide.md)
- [EC2 Deployment Guide](../documentation/ec2-deployment-guide.md)
- [Unity Marketplace](https://github.com/unity-sds/unity-marketplace)
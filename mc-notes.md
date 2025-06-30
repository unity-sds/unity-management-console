# S3 Backend for Terraform State Files

The Unity Management Console uses an AWS S3 bucket as the backend for storing Terraform state files, with DynamoDB for state locking. This is a common best practice for managing Terraform state in production environments. Here are the specifics:

### 1. State File Storage
- S3 Bucket: State files are stored in an S3 bucket that's created during the bootstrap process
- Bucket Naming: The bucket name is determined in one of two ways:
  - Provided explicitly through the `BucketName` configuration
  - Retrieved from an SSM parameter at path `/unity/{project}/{venue}/cs/monitoring/s3/bucketName`
  - If not found, the system previously used a randomized name with the prefix `mgmt-` followed by 8 random alphanumeric characters

### 2. State File Structure
- Each component's state is stored with a key using the format: `{project}-{venue}-tfstate`
- This ensures isolation between different environments (venues) and projects

### 3. State Locking with DynamoDB
- A DynamoDB table is created to handle state locking
- The table name follows the pattern: `{project}-{venue}-terraform-state`
- This prevents concurrent modifications to the same state which could cause conflicts

### 4. Configuration Setup
In the `writeInitTemplate` function in `bootstrap.go`, the Terraform backend configuration is set up as follows:

The S3 bucket name, key format, and AWS region are passed as runtime parameters to Terraform when initializing:

Where:
- bucket is `bucket={bucketName}`
- key is `key={project}-{venue}-tfstate`
- region is `region={awsRegion}`

### 5. Security and Versioning
The S3 bucket is configured with:
- **TLS-only access**: A bucket policy ensures that only secure TLS connections are allowed
- **Versioning**: S3 versioning is enabled to track all state file changes
- **Lifecycle policies**: For certain prefixes like health checks, lifecycle policies are set to manage retention

### 6. Per-Application State Management
For individual applications installed via the management console:

- The `AddApplicationToStack` function creates Terraform configurations in separate files
- Each application gets a file named `{applicationName}-{deploymentName}.tf`
- These all use the same S3 backend, with the components distinguished by module names in the state

This architecture ensures that the state for each component's infrastructure is securely stored, versioned, and protected from concurrent modifications, while still allowing the Unity Management Console to manage everything centrally.
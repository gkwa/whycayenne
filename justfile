# Justfile for Terraform DynamoDB project
# Use `just format-all` to format this file and Terraform files

# Default recipe to run when just is called without arguments
default:
    @echo "Available recipes:"
    @just --list

# Format this Justfile
format-justfile:
    just --unstable --fmt

# Format Terraform files
format-terraform:
    terraform fmt -recursive

# Format all files (Justfile and Terraform)
format-all: format-justfile format-terraform
    @echo "All files formatted."

# Initialize Terraform
init:
    terraform init

# Validate Terraform files
validate:
    terraform validate

# Plan Terraform changes and save to tfplan
plan:
    terraform plan -out=tfplan

# Apply Terraform changes
apply:
    terraform apply tfplan

# Destroy Terraform resources
destroy:
    terraform destroy

# Show Terraform state
show:
    terraform show

# List Terraform resources
list:
    terraform state list

# Clean up Terraform files
clean:
    rm -rf .terraform tfplan

# Full cleanup: remove all generated files and Terraform state
cleanup:
    rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup tfplan
    @echo "Cleanup complete. Project reset to initial state."

# Run all pre-apply checks
check: format-all validate plan

# Full cycle: check, apply, and show
cycle: check apply show

# Get DynamoDB table info (adjust the table name if necessary)
table-info:
    aws dynamodb describe-table --table-name ProductTable

# Estimate cost (requires infracost to be installed)
cost:
    infracost breakdown --path .

# Run security scan (requires tfsec to be installed)
scan:
    tfsec .

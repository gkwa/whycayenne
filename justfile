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
    terraform -chdir=terraform fmt

# Format all files (Justfile and Terraform)
format-all: format-justfile format-terraform
    @echo "All files formatted."

# Initialize Terraform
init:
    terraform -chdir=terraform init

# Validate Terraform files
validate:
    terraform -chdir=terraform validate

# Plan Terraform changes and save to tfplan
plan:
    terraform -chdir=terraform plan -out=tfplan

# Apply Terraform changes from saved plan
apply:
    terraform -chdir=terraform apply tfplan

# Plan and apply Terraform changes
plan-apply: plan apply

# Destroy Terraform resources (interactive)
destroy:
    terraform -chdir=terraform destroy

# Plan destroy and save to destroy.tfplan
plan-destroy:
    terraform -chdir=terraform plan -destroy -out=destroy.tfplan

# Apply destroy plan without prompting
destroy-auto: plan-destroy
    terraform -chdir=terraform apply -auto-approve destroy.tfplan

# Show Terraform state
show:
    terraform -chdir=terraform show

# List Terraform resources
list:
    terraform -chdir=terraform state list

# Clean up Terraform files
clean:
    rm -rf terraform/.terraform terraform/tfplan terraform/destroy.tfplan

# Full cleanup: remove all generated files and Terraform state
cleanup:
    rm -rf terraform/.terraform terraform/.terraform.lock.hcl terraform/terraform.tfstate terraform/terraform.tfstate.backup terraform/tfplan terraform/destroy.tfplan
    @echo "Cleanup complete. Project reset to initial state."

# Run all pre-apply checks
check: format-all validate plan

# Full cycle: check, plan-apply, and show
cycle: check plan-apply show

# Get DynamoDB table info (adjust the table name if necessary)
table-info:
    aws dynamodb describe-table --table-name ProductTable

# Estimate cost (requires infracost to be installed)
cost:
    infracost breakdown --path terraform

# Run security scan (requires tfsec to be installed)
scan:
    tfsec terraform

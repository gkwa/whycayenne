# Configure the AWS provider
provider "aws" {
  region = "ca-central-1"
}

# Create the DynamoDB table
resource "aws_dynamodb_table" "product_table" {
  name           = "ProductTable"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "GSI1PK"
    type = "S"
  }

  attribute {
    name = "GSI1SK"
    type = "S"
  }

  # Global Secondary Index for category-based queries
  global_secondary_index {
    name               = "GSI1"
    hash_key           = "GSI1PK"
    range_key          = "GSI1SK"
    projection_type    = "ALL"
  }

  tags = {
    Name        = "product-table"
    Environment = "Production"
  }
}

# Output the table name
output "table_name" {
  value       = aws_dynamodb_table.product_table.name
  description = "Name of the DynamoDB table"
}

# Output the table ARN
output "table_arn" {
  value       = aws_dynamodb_table.product_table.arn
  description = "ARN of the DynamoDB table"
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ca-central-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Define the query input
	input := &dynamodb.QueryInput{
		TableName:              aws.String("ProductTable"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :category"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category": &types.AttributeValueMemberS{Value: "CATEGORY#pepper"},
		},
	}

	// Execute the query
	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error querying table: %s", err)
	}

	// Process the results
	fmt.Printf("Query returned %d items.\n", len(result.Items))
	for _, item := range result.Items {
		name := item["name"].(*types.AttributeValueMemberS).Value
		price := item["price"].(*types.AttributeValueMemberN).Value
		store := item["store"].(*types.AttributeValueMemberS).Value
		fmt.Printf("Name: %s, Price: $%s, Store: %s\n", name, price, store)
	}
}
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func queryData(svc *dynamodb.Client, queryString string) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("ProductTable"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :category"),
		FilterExpression:       aws.String("contains(#name_lower, :queryString)"),
		ExpressionAttributeNames: map[string]string{
			"#name_lower": "name_lower",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category":    &types.AttributeValueMemberS{Value: "CATEGORY#pepper"},
			":queryString": &types.AttributeValueMemberS{Value: strings.ToLower(queryString)},
		},
	}

	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error querying table: %s", err)
	}

	totalCount, err := getTotalCount(svc)
	if err != nil {
		log.Printf("Error getting total count: %s", err)
		totalCount = 0
	}
	fmt.Printf("Query returned %d items out of %d total records.\n", len(result.Items), totalCount)
	for _, item := range result.Items {
		printItem(item)
	}
}

func queryAllPeppers(svc *dynamodb.Client) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("ProductTable"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :category"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category": &types.AttributeValueMemberS{Value: "CATEGORY#pepper"},
		},
	}

	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error querying table: %s", err)
	}

	totalCount, err := getTotalCount(svc)
	if err != nil {
		log.Printf("Error getting total count: %s", err)
		totalCount = 0
	}
	fmt.Printf("Query returned %d items out of %d total records.\n", len(result.Items), totalCount)
	for _, item := range result.Items {
		printItem(item)
	}
}

func printItem(item map[string]types.AttributeValue) {
	name := item["name"].(*types.AttributeValueMemberS).Value
	price := item["price"].(*types.AttributeValueMemberN).Value
	store := item["store"].(*types.AttributeValueMemberS).Value
	fmt.Printf("Name: %s, Price: $%s, Store: %s\n", name, price, store)
}

func getTotalCount(svc *dynamodb.Client) (int64, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("ProductTable"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "METADATA"},
			"SK": &types.AttributeValueMemberS{Value: "COUNTER"},
		},
	}

	result, err := svc.GetItem(context.TODO(), input)
	if err != nil {
		return 0, err
	}

	if count, ok := result.Item["count"].(*types.AttributeValueMemberN); ok {
		return strconv.ParseInt(count.Value, 10, 64)
	}

	return 0, fmt.Errorf("count not found")
}

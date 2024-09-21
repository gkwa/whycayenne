package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func unloadData(svc *dynamodb.Client) {
	var lastEvaluatedKey map[string]types.AttributeValue
	totalDeleted := 0

	for {
		scanInput := &dynamodb.ScanInput{
			TableName:         aws.String("ProductTable"),
			Limit:             aws.Int32(25), // Scan up to 25 items at a time
			ExclusiveStartKey: lastEvaluatedKey,
		}

		result, err := svc.Scan(context.TODO(), scanInput)
		if err != nil {
			log.Fatalf("Got error scanning table: %s", err)
		}

		if len(result.Items) == 0 {
			break
		}

		deleteRequests := make([]types.WriteRequest, 0, len(result.Items))
		for _, item := range result.Items {
			deleteRequests = append(deleteRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"PK": item["PK"],
						"SK": item["SK"],
					},
				},
			})
		}

		batchWriteInput := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				"ProductTable": deleteRequests,
			},
		}

		_, err = svc.BatchWriteItem(context.TODO(), batchWriteInput)
		if err != nil {
			log.Printf("Got error calling BatchWriteItem: %s", err)
		}

		totalDeleted += len(result.Items)
		fmt.Printf("Deleted %d items\n", totalDeleted)

		lastEvaluatedKey = result.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}
	}

	// Update the counter
	updateCounterInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("ProductTable"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "METADATA"},
			"SK": &types.AttributeValueMemberS{Value: "COUNTER"},
		},
		UpdateExpression: aws.String("SET #count = :zero"),
		ExpressionAttributeNames: map[string]string{
			"#count": "count",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":zero": &types.AttributeValueMemberN{Value: "0"},
		},
	}

	_, err := svc.UpdateItem(context.TODO(), updateCounterInput)
	if err != nil {
		log.Printf("Got error updating counter: %s", err)
	}

	fmt.Printf("All data removed from the table. Total items deleted: %d\n", totalDeleted)
}

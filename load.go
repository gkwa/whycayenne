package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func loadData(svc *dynamodb.Client, verbose bool) {
	file, err := os.Open("data.jsonl")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var product Product
		err := json.Unmarshal(scanner.Bytes(), &product)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		input := &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Put: &types.Put{
						TableName: aws.String("ProductTable"),
						Item:      productToItem(product),
					},
				},
				{
					Update: &types.Update{
						TableName: aws.String("ProductTable"),
						Key: map[string]types.AttributeValue{
							"PK": &types.AttributeValueMemberS{Value: "METADATA"},
							"SK": &types.AttributeValueMemberS{Value: "COUNTER"},
						},
						UpdateExpression: aws.String("ADD #count :inc"),
						ExpressionAttributeNames: map[string]string{
							"#count": "count",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":inc": &types.AttributeValueMemberN{Value: "1"},
						},
					},
				},
			},
		}

		_, err = svc.TransactWriteItems(context.TODO(), input)
		if err != nil {
			log.Printf("Got error calling TransactWriteItems: %s", err)
		}

		if verbose {
			log.Printf("Inserted item: %s", product.Name)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data loaded successfully")
}

func loadDataBatch(svc *dynamodb.Client, verbose bool) {
	file, err := os.Open("data.jsonl")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	batch := make([]types.WriteRequest, 0, 25)
	uniqueKeys := make(map[string]bool)
	totalItems := 0

	for scanner.Scan() {
		var product Product
		err := json.Unmarshal(scanner.Bytes(), &product)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		key := fmt.Sprintf("%s#%s", product.Name, product.Store)
		if uniqueKeys[key] {
			log.Printf("Duplicate item found: %s. Skipping.", key)
			continue
		}
		uniqueKeys[key] = true

		batch = append(batch, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: productToItem(product),
			},
		})

		if len(batch) == 25 {
			if err := writeBatch(svc, batch, verbose); err != nil {
				log.Printf("Error writing batch: %v", err)
			}
			totalItems += len(batch)
			batch = batch[:0]
			uniqueKeys = make(map[string]bool)
		}
	}

	if len(batch) > 0 {
		if err := writeBatch(svc, batch, verbose); err != nil {
			log.Printf("Error writing final batch: %v", err)
		}
		totalItems += len(batch)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Update the counter
	if err := updateCounter(svc, totalItems); err != nil {
		log.Printf("Error updating counter: %v", err)
	}

	fmt.Printf("Data loaded successfully. Total items: %d\n", totalItems)
}

func writeBatch(svc *dynamodb.Client, batch []types.WriteRequest, verbose bool) error {
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"ProductTable": batch,
		},
	}

	_, err := svc.BatchWriteItem(context.TODO(), input)
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Batch write: %d items", len(batch))
	}

	return nil
}

func updateCounter(svc *dynamodb.Client, count int) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("ProductTable"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "METADATA"},
			"SK": &types.AttributeValueMemberS{Value: "COUNTER"},
		},
		UpdateExpression: aws.String("ADD #count :inc"),
		ExpressionAttributeNames: map[string]string{
			"#count": "count",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", count)},
		},
	}

	_, err := svc.UpdateItem(context.TODO(), input)
	return err
}

func productToItem(product Product) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%s", product.Name)},
		"SK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("METADATA#%s", product.Store)},
		"GSI1PK":       &types.AttributeValueMemberS{Value: "CATEGORY#pepper"},
		"GSI1SK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%s", product.Name)},
		"name":         &types.AttributeValueMemberS{Value: product.Name},
		"name_lower":   &types.AttributeValueMemberS{Value: strings.ToLower(product.Name)},
		"price":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.Price)},
		"price_per_lb": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.PricePerLb)},
		"price_per_oz": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.PricePerOz)},
		"store":        &types.AttributeValueMemberS{Value: product.Store},
		"volume":       &types.AttributeValueMemberS{Value: product.Volume},
		"weight":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.Weight)},
		"on_sale":      &types.AttributeValueMemberBOOL{Value: product.OnSale},
		"datetime":     &types.AttributeValueMemberS{Value: product.DateTime},
	}
}

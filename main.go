package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Product struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	PricePerLb float64 `json:"price_per_lb"`
	PricePerOz float64 `json:"price_per_oz"`
	Store      string  `json:"store"`
	Volume     string  `json:"volume"`
	Weight     float64 `json:"weight"`
	OnSale     bool    `json:"on_sale"`
	DateTime   string  `json:"datetime"`
}

func main() {
	loadFlag := flag.Bool("load", false, "Load data into DynamoDB")
	unloadFlag := flag.Bool("unload", false, "Remove all data from DynamoDB")
	queryFlag := flag.String("query", "", "Query by pepper name")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ca-central-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	if *loadFlag {
		loadData(svc)
	} else if *unloadFlag {
		unloadData(svc)
	} else if *queryFlag != "" {
		queryData(svc, *queryFlag)
	} else {
		queryAllPeppers(svc)
	}
}

func loadData(svc *dynamodb.Client) {
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

		input := &dynamodb.PutItemInput{
			TableName: aws.String("ProductTable"),
			Item: map[string]types.AttributeValue{
				"PK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%s", product.Name)},
				"SK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("METADATA#%s", product.Store)},
				"GSI1PK":       &types.AttributeValueMemberS{Value: "CATEGORY#pepper"},
				"GSI1SK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%s", product.Name)},
				"name":         &types.AttributeValueMemberS{Value: product.Name},
				"price":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.Price)},
				"price_per_lb": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.PricePerLb)},
				"price_per_oz": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.PricePerOz)},
				"store":        &types.AttributeValueMemberS{Value: product.Store},
				"volume":       &types.AttributeValueMemberS{Value: product.Volume},
				"weight":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.Weight)},
				"on_sale":      &types.AttributeValueMemberBOOL{Value: product.OnSale},
				"datetime":     &types.AttributeValueMemberS{Value: product.DateTime},
			},
		}

		_, err = svc.PutItem(context.TODO(), input)
		if err != nil {
			log.Printf("Got error calling PutItem: %s", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data loaded successfully")
}

func unloadData(svc *dynamodb.Client) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("ProductTable"),
	}

	result, err := svc.Scan(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error scanning table: %s", err)
	}

	for _, item := range result.Items {
		deleteInput := &dynamodb.DeleteItemInput{
			TableName: aws.String("ProductTable"),
			Key: map[string]types.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			},
		}

		_, err := svc.DeleteItem(context.TODO(), deleteInput)
		if err != nil {
			log.Printf("Got error calling DeleteItem: %s", err)
		}
	}

	fmt.Println("All data removed from the table")
}

func queryData(svc *dynamodb.Client, queryString string) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("ProductTable"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :category AND begins_with(GSI1SK, :productPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category":      &types.AttributeValueMemberS{Value: "CATEGORY#pepper"},
			":productPrefix": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%s", queryString)},
		},
	}

	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error querying table: %s", err)
	}

	fmt.Printf("Query returned %d items.\n", len(result.Items))
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

	fmt.Printf("Query returned %d items.\n", len(result.Items))
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

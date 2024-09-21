package main

import (
	"context"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	loadFlag := flag.Bool("load", false, "Load data into DynamoDB")
	loadBatchFlag := flag.Bool("loadbatch", false, "Load data into DynamoDB using batch write")
	unloadFlag := flag.Bool("unload", false, "Remove all data from DynamoDB")
	queryFlag := flag.String("query", "", "Query by pepper name")
	verboseFlag := flag.Bool("v", false, "Enable verbose logging")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ca-central-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	if *loadFlag {
		loadData(svc, *verboseFlag)
	} else if *loadBatchFlag {
		loadDataBatch(svc, *verboseFlag)
	} else if *unloadFlag {
		unloadData(svc)
	} else if *queryFlag != "" {
		queryData(svc, *queryFlag)
	} else {
		queryAllPeppers(svc)
	}
}

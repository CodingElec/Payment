package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB struct {
	Cli       *dynamodb.Client
	TableName string
	IdempoTbl string
}

func New(ctx context.Context, table, idempo string) (*DB, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}

	endpoint := "http://localhost:4566" //hardcoded for now. OS env or smth else later.
	clientOpts := func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	}

	return &DB{
		Cli:       dynamodb.NewFromConfig(cfg, clientOpts),
		TableName: table,
		IdempoTbl: idempo,
	}, nil
}

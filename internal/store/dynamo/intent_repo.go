package dynamo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/raf-fml/Payment/internal/domain"
)

func (db *DB) PutPaymentIntent(ctx context.Context, merchantID string, req domain.PaymentIntent) (domain.PaymentIntent, error) {
	now := time.Now().UTC()
	intent := domain.PaymentIntent{
		IntentID:    "pi_" + uuid.NewString(),
		MerchantID:  merchantID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		Status:      domain.IntentCreated,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	pk := fmt.Sprintf("PINT#%s", intent.IntentID)

	item, _ := attributevalue.MarshalMap(struct {
		PK   string `dynamodbav:"pk"`
		SK   string `dynamodbav:"sk"`
		Type string `dynamodbav:"type"`
		domain.PaymentIntent
	}{
		PK:            pk,
		SK:            "META",
		Type:          "PaymentIntent",
		PaymentIntent: intent,
	})

	ev := domain.OutboxEvent{
		EventID:  "ev_" + uuid.NewString(),
		Type:     domain.EvPaymentIntentCreated,
		IntentID: intent.IntentID,
	}
	payload, _ := json.Marshal(intent)
	ev.Payload = payload
	evItem := marshalOutboxItem(intent.IntentID, ev) //need to create this on outbox repo.

	_, err := db.Cli.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Put: &types.Put{
				TableName:           &db.TableName,
				Item:                item,
				ConditionExpression: awsString("attribute_not_exists(pk)"),
			}},
			{Put: &types.Put{TableName: &db.TableName, Item: evItem}},
		},
	})
	return intent, err
}

func (db *DB) GetPaymentIntent(ctx context.Context, intentID string) (domain.PaymentIntent, error) {
	out, err := db.Cli.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &db.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PINT#%s", intentID)},
			"sk": &types.AttributeValueMemberS{Value: "META"},
		},
	})
	if err != nil || out.Item == nil {
		return domain.PaymentIntent{}, fmt.Errorf("not found")
	}
	var row struct {
		Type string `dynamodbav:"type"`
		domain.PaymentIntent
	}
	if err := attributevalue.UnmarshalMap(out.Item, &row); err != nil {
		return domain.PaymentIntent{}, err
	}
	return row.PaymentIntent, nil
}

func awsString(s string) *string { return &s }

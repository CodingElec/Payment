package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/raf-fml/Payment/internal/domain"
)

func marshalOutboxItem(intentID string, ev domain.OutboxEvent) map[string]types.AttributeValue {
	pk := "OUTBOX#" + intentID
	item, _ := attributevalue.MarshalMap(struct {
		PK     string `dynamodbav:"pk"`
		SK     string `dynamodbav:"sk"`
		Type   string `dynamodbav:"type"`
		Status string `dynamodbav:"status"`
		domain.OutboxEvent
	}{
		PK:          pk,
		SK:          ev.EventID,
		Type:        "Outbox",
		Status:      "pending",
		OutboxEvent: ev,
	})
	return item
}

func (db *DB) FetchPendingOutbox(ctx context.Context, limit int32) ([]domain.OutboxEvent, error) {
	out, err := db.Cli.Scan(ctx, &dynamodb.ScanInput{
		TableName:                &db.TableName,
		FilterExpression:         awsString("#t = :obox AND #st = :pending"),
		ExpressionAttributeNames: map[string]string{"#t": "type", "#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":obox":    &types.AttributeValueMemberS{Value: "Outbox"},
			":pending": &types.AttributeValueMemberS{Value: "pending"},
		},
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	var evs []domain.OutboxEvent
	for _, m := range out.Items {
		var row struct {
			domain.OutboxEvent `dynamodbav:"OutboxEvent,inline"`
			EventID            string `dynamodbav:"eventId"`
			IntentID           string `dynamodbav:"intentId"`
			Type2              string `dynamodbav:"type2"`
			Payload            []byte `dynamodbav:"payload"`
		}
		_ = attributevalue.UnmarshalMap(m, &row)

		// Some drivers require manual extraction:
		if row.EventID == "" {
			if v, ok := m["eventId"].(*types.AttributeValueMemberS); ok {
				row.EventID = v.Value
			}
		}
		if row.IntentID == "" {
			if v, ok := m["intentId"].(*types.AttributeValueMemberS); ok {
				row.IntentID = v.Value
			}
		}
		if len(row.Payload) == 0 {
			if v, ok := m["payload"].(*types.AttributeValueMemberB); ok {
				row.Payload = v.Value
			}
		}
		typ := row.Type
		if typ == "" && row.Type2 != "" {
			typ = domain.EventType(row.Type2)
		}
		evs = append(evs, domain.OutboxEvent{
			EventID:  row.EventID,
			Type:     typ,
			IntentID: row.IntentID,
			Payload:  row.Payload,
		})
	}
	return evs, nil
}

func (db *DB) MarkOutboxShipped(ctx context.Context, intentID, eventID string) error {
	_, err := db.Cli.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &db.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "OUTBOX#" + intentID},
			"sk": &types.AttributeValueMemberS{Value: eventID},
		},
		UpdateExpression:         awsString("SET #st = :done"),
		ExpressionAttributeNames: map[string]string{"#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":done": &types.AttributeValueMemberS{Value: "shipped"},
		},
	})
	return err
}

// helper puts for transact
func evPut(tbl string, item map[string]types.AttributeValue) *types.Put {
	return &types.Put{TableName: &tbl, Item: item}
}

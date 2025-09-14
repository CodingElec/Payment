package domain

type EventType string

const (
	EvPaymentIntentCreated EventType = "PaymentIntentCreated"
	EvTransactionCreated   EventType = "TransactionCreated"
	EvTransactionCompleted EventType = "TransactionCompleted"
)

type OutboxEvent struct {
	EventID  string    `json:"eventId"`
	Type     EventType `json:"type"`
	IntentID string    `json:"intentId"`
	Payload  []byte    `json:"payload"`
}

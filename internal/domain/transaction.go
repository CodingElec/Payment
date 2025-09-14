package domain

import "time"

type TransactionStatus string

const (
	TxPending   TransactionStatus = "pending"
	TxSucceeded TransactionStatus = "succeeded"
	TxFailed    TransactionStatus = "failed"
)

type Transaction struct {
	TransactionID string            `json:"transactionId"`
	IntentID      string            `json:"intentId"`
	Type          string            `json:"type"`
	Status        TransactionStatus `json:"status"`
	Amount        int64             `json:"amountInCents"`
	Currency      string            `json:"currency"`
	ErrorMessage  string            `json:"errorMessage,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

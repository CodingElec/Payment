package domain

import "time"

type PaymentIntentStatus string

const (
	IntentCreated    PaymentIntentStatus = "created"
	IntentProcessing PaymentIntentStatus = "processed"
	IntentSucceeded  PaymentIntentStatus = "succeeded"
	IntentFaile      PaymentIntentStatus = "failed"
)

type PaymentIntent struct {
	IntentID    string              `json:"paymentIntentId"`
	MerchantID  string              `json:"merchantId"`
	Amount      int64               `json:"amountInCents"`
	Currency    string              `json:"currency"`
	Description string              `json:"description,omitempty"`
	Status      PaymentIntentStatus `json:"status"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

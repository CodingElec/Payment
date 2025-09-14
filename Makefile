.PHONY: run-intent run-txn run-outbox create-tables list-tables

AWS_ENDPOINT=http://localhost:4566
REGION=us-east-1

run-intent:
	go run ./cmd/paymentintent

run-txn:
	go run ./cmd/transaction

run-outbox:
	go run ./cmd/outbox

create-tables:
	AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test aws --endpoint-url=$(AWS_ENDPOINT) dynamodb create-table \
	  --region $(REGION) --table-name PaymentTable \
	  --attribute-definitions AttributeName=pk,AttributeType=S AttributeName=sk,AttributeType=S \
	  --key-schema AttributeName=pk,KeyType=HASH AttributeName=sk,KeyType=RANGE \
	  --billing-mode PAY_PER_REQUEST || true

	AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test aws --endpoint-url=$(AWS_ENDPOINT) dynamodb create-table \
	  --region $(REGION) --table-name IdempotencyTable \
	  --attribute-definitions AttributeName=pk,AttributeType=S \
	  --key-schema AttributeName=pk,KeyType=HASH \
	  --billing-mode PAY_PER_REQUEST || true

list-tables:
	AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test aws --endpoint-url=$(AWS_ENDPOINT) --region $(REGION) dynamodb list-tables

# Load optional .env (dev-friendly)
-include .env
export

.PHONY: run-intent run-txn run-outbox create-tables list-tables


AWS_DDB = aws --endpoint-url=$(AWS_ENDPOINT) --region $(REGION) dynamodb

run-intent:
	go run ./cmd/paymentintent

run-txn:
	go run ./cmd/transaction

run-outbox:
	go run ./cmd/outbox

create-tables:
	-$(AWS_DDB) create-table --table-name PaymentTable --attribute-definitions AttributeName=pk,AttributeType=S AttributeName=sk,AttributeType=S --key-schema AttributeName=pk,KeyType=HASH AttributeName=sk,KeyType=RANGE --billing-mode PAY_PER_REQUEST
	-$(AWS_DDB) create-table --table-name IdempotencyTable --attribute-definitions AttributeName=pk,AttributeType=S --key-schema AttributeName=pk,KeyType=HASH --billing-mode PAY_PER_REQUEST

list-tables:
	$(AWS_DDB) list-tables

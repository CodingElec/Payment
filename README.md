# Payment
Payment Service Sandbox

# Objective
Implement a Payment Service that:
## Core Requirements 

Merchants should be able to initiate payment requests (charge a customer for a specific amount). 
Users should be able to pay for products with credit/debit cards. Merchants should be able to view status updates for payments (e.g., pending, success, failed). 
## Below the line (out of scope): 
Customers should be able to save payment methods for future use. Merchants should be able to issue full or partial refunds. Merchants should be able to view transaction history and reports. Support for alternative payment methods (e.g., bank transfers, digital wallets). Handling recurring payments (subscriptions). Payouts to merchants.

# Payment System Flow (End-to-End) - Local Test / Sandbox

```mermaid
sequenceDiagram
autonumber
participant Client
participant PI as PaymentIntentSvc
participant TX as TransactionSvc
participant IDEM as IdempotencyTable
participant DB as PaymentTable
participant SHIP as OutboxShipper
participant KAFKA as Kafka
participant PROC as Processor

Client->>PI: POST /payment-intents
PI->>DB: Put PaymentIntent + Outbox
PI-->>Client: 201 created

Client->>TX: POST /payment-intents/:id/transactions (Idempotency-Key)
TX->>IDEM: Check key
IDEM-->>TX: Not found
TX->>DB: Put Transaction + Update Intent + Outbox
TX-->>Client: 202 pending
TX->>IDEM: Store response

Client->>TX: POST again with same key
TX->>IDEM: Lookup
IDEM-->>TX: Hit
TX-->>Client: 200 replayed

SHIP->>DB: Scan Outbox (pending)
DB-->>SHIP: Events
SHIP->>KAFKA: Produce event
SHIP->>DB: Mark shipped

PROC->>DB: Scan pending Txns
DB-->>PROC: Results
PROC->>DB: Update txn -> succeeded
PROC->>DB: Update intent -> succeeded

Client->>PI: GET /payment-intents/:id
PI->>DB: Get PaymentIntent
DB-->>PI: succeeded
PI-->>Client: 200 JSON
```

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/you/payments/internal/domain"
	"github.com/you/payments/internal/store/dynamo"
)

type createIntentReq struct {
	Amount      int64  `json:"amountInCents" binding:"required,gt=0"`
	Currency    string `json:"currency" binding:"required"`
	Description string `json:"description"`
}

func main() {
	ctx := context.Background()
	db, err := dynamo.New(ctx, "PaymentTable", "IdempotencyTable")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery(), gin.LoggerWithWriter(gin.DefaultWriter, "/healthz"))
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.POST("/payment-intents", func(c *gin.Context) {
		var in createIntentReq
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		intent, err := db.PutPaymentIntent(c.Request.Context(), "merch_demo",
			domain.PaymentIntent{Amount: in.Amount, Currency: in.Currency, Description: in.Description})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"paymentIntentId": intent.IntentID,
			"status":          intent.Status,
			"createdAt":       intent.CreatedAt.Format(time.RFC3339),
		})
	})

	r.GET("/payment-intents/:id", func(c *gin.Context) {
		id := c.Param("id")
		intent, err := db.GetPaymentIntent(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, intent)
	})

	log.Println("paymentintent service on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

package payment

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"kitchen-api/internal/app/order"
)

type PaymentController struct {
	paymentService *PaymentService
}

func NewPaymentController(paymentService *PaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

func (c *PaymentController) CreateCheckoutSession(ctx *gin.Context) {
	var req order.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.paymentService.CreateCheckoutSession(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Checkout session created successfully",
		"data":    result,
		"success": true,
	})
}

func (c *PaymentController) HandleWebhook(ctx *gin.Context) {
	const maxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxBodyBytes)

	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("[Stripe Webhook] Error reading body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	sigHeader := ctx.GetHeader("Stripe-Signature")
	log.Printf("[Stripe Webhook] Received event, payload size: %d, sig present: %v", len(payload), sigHeader != "")
	if err := c.paymentService.HandleWebhook(payload, sigHeader); err != nil {
		log.Printf("[Stripe Webhook] Error processing event: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Stripe Webhook] Successfully processed and confirmed order")
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

func (c *PaymentController) GetSessionStatus(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	status, err := c.paymentService.GetSessionStatus(sessionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    status,
		"success": true,
		"message": "Payment session fetched successfully",
	})
}

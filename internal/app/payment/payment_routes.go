package payment

import (
	"github.com/gin-gonic/gin"
	"kitchen-api/internal/app/order"
	"kitchen-api/internal/database"
)

func RegisterPaymentRoutes(r *gin.Engine, db *database.OrmDb, orderService *order.OrderService) {
	paymentRepo := NewPaymentRepository(db)
	service := NewPaymentService(paymentRepo, orderService)
	controller := NewPaymentController(service)

	paymentGroup := r.Group("/payments")
	{
		paymentGroup.POST("/checkout", controller.CreateCheckoutSession)
		paymentGroup.POST("/webhook", controller.HandleWebhook)
		paymentGroup.GET("/session/:sessionId", controller.GetSessionStatus)
	}
}

package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *OrderService
}

func NewOrderController(orderService *OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.orderService.CreateOrder(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Order placed successfully",
		"data":    result,
		"success": true,
	})
}

func (c *OrderController) GetOrdersByRestaurant(ctx *gin.Context) {
	restaurantID := ctx.Param("id")
	if restaurantID == "" {
		restaurantID = ctx.Param("restaurantId")
	}
	orders, err := c.orderService.GetOrdersByRestaurantID(restaurantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    orders,
		"success": true,
		"message": "Orders fetched successfully",
	})
}

func (c *OrderController) GetOrdersByTable(ctx *gin.Context) {
	tableID := ctx.Param("id")
	if tableID == "" {
		tableID = ctx.Param("tableId")
	}
	orders, err := c.orderService.GetOrdersByTableID(tableID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    orders,
		"success": true,
		"message": "Orders fetched successfully",
	})
}

func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	o, err := c.orderService.GetOrderByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    o,
		"success": true,
		"message": "Order fetched successfully",
	})
}

func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.orderService.UpdateOrderStatus(id, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"success": true,
	})
}

package orderitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderItemController struct {
	service *OrderItemService
}

func NewOrderItemController(service *OrderItemService) *OrderItemController {
	return &OrderItemController{service: service}
}

func (c *OrderItemController) GetItemsByOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")
	if orderID == "" {
		orderID = ctx.Param("orderId")
	}
	items, err := c.service.GetItemsByOrderID(orderID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    items,
		"success": true,
		"message": "Order items fetched successfully",
	})
}

func (c *OrderItemController) GetOrderItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := c.service.GetOrderItemByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    item,
		"success": true,
		"message": "Order item fetched successfully",
	})
}

package order

import "kitchen-api/internal/enums"

type OrderItemAddonInput struct {
	AddonID  string `json:"addonId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type OrderItemInput struct {
	MenuItemID string                `json:"menuItemId" binding:"required"`
	Quantity   int                   `json:"quantity" binding:"required,gt=0"`
	Addons     []OrderItemAddonInput `json:"addons,omitempty"`
}

type CreateOrderRequest struct {
	RestaurantID string           `json:"restaurantId" binding:"required"`
	TableID      string           `json:"tableId" binding:"required"`
	Notes        string           `json:"notes,omitempty"`
	Items        []OrderItemInput `json:"items" binding:"required,min=1,dive"`
}

type UpdateOrderStatusRequest struct {
	Status enums.OrderStatus `json:"status" binding:"required"`
}

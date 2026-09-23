package orderitem

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderItemService struct {
	repo *OrderItemRepository
}

func NewOrderItemService(repo *OrderItemRepository) *OrderItemService {
	return &OrderItemService{repo: repo}
}

func (s *OrderItemService) GetItemsByOrderID(orderIDStr string) ([]OrderItem, error) {
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID")
	}
	return s.repo.GetItemsByOrderID(orderID)
}

func (s *OrderItemService) GetOrderItemByID(idStr string) (*OrderItem, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order item ID")
	}
	item, err := s.repo.GetOrderItemByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("order item not found")
	}
	return item, nil
}

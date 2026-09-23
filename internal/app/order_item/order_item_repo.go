package orderitem

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ktmbeestech/yanshi/database"
	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *database.OrmDb
}

func NewOrderItemRepository(db *database.OrmDb) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (repo *OrderItemRepository) CreateOrderItem(item *OrderItem) error {
	return repo.db.OrmInstance.Create(item).Error
}

func (repo *OrderItemRepository) CreateOrderItemsBatch(items []OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return repo.db.OrmInstance.Create(&items).Error
}

func (repo *OrderItemRepository) GetItemsByOrderID(orderID uuid.UUID) ([]OrderItem, error) {
	var items []OrderItem
	err := repo.db.OrmInstance.
		Preload("MenuItem").
		Preload("Addons.Addon").
		Where("order_id = ?", orderID).
		Find(&items).Error
	return items, err
}

func (repo *OrderItemRepository) GetOrderItemByID(id uuid.UUID) (*OrderItem, error) {
	var item OrderItem
	err := repo.db.OrmInstance.
		Preload("MenuItem").
		Preload("Addons.Addon").
		First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

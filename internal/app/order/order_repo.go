package order

import (
	"errors"

	orderitem "kitchen-api/internal/app/order_item"
	"kitchen-api/internal/enums"

	"github.com/google/uuid"
	"github.com/ktmbeestech/yanshi/database"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *database.OrmDb
}

func NewOrderRepository(db *database.OrmDb) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) CreateOrderWithItems(o *Order, items []orderitem.OrderItem) error {
	return repo.db.OrmInstance.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = o.ID
			if items[i].ID == uuid.Nil {
				items[i].ID = uuid.New()
			}
			for j := range items[i].Addons {
				items[i].Addons[j].OrderItemID = items[i].ID
				if items[i].Addons[j].ID == uuid.Nil {
					items[i].Addons[j].ID = uuid.New()
				}
			}
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		o.OrderItems = items
		return nil
	})
}

func (repo *OrderRepository) GetOrdersByRestaurantID(restaurantID uuid.UUID) ([]Order, error) {
	var orders []Order
	err := repo.db.OrmInstance.
		Preload("OrderItems.MenuItem").
		Preload("OrderItems.Addons.Addon").
		Preload("Table").
		Preload("Restaurant").
		Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (repo *OrderRepository) GetOrdersByTableID(tableID uuid.UUID) ([]Order, error) {
	var orders []Order
	err := repo.db.OrmInstance.
		Preload("OrderItems.MenuItem").
		Preload("OrderItems.Addons.Addon").
		Preload("Table").
		Preload("Restaurant").
		Where("table_id = ?", tableID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (repo *OrderRepository) GetOrderByID(id uuid.UUID) (*Order, error) {
	var o Order
	err := repo.db.OrmInstance.
		Preload("OrderItems.MenuItem").
		Preload("OrderItems.Addons.Addon").
		Preload("Table").
		Preload("Restaurant").
		First(&o, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (repo *OrderRepository) UpdateOrderStatus(id uuid.UUID, status enums.OrderStatus) error {
	return repo.db.OrmInstance.Model(&Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

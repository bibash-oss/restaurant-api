package menuitem

import (
	"errors"

	"github.com/google/uuid"
	"kitchen-api/internal/database"
	"gorm.io/gorm"
)

type MenuItemRepository struct {
	db *database.OrmDb
}

func NewMenuItemRepository(db *database.OrmDb) *MenuItemRepository {
	return &MenuItemRepository{db: db}
}

func (repo *MenuItemRepository) CreateMenuItem(item *MenuItem) error {
	return repo.db.OrmInstance.Create(item).Error
}

func (repo *MenuItemRepository) GetMenuItemsByRestaurantID(restaurantID uuid.UUID) ([]MenuItem, error) {
	var items []MenuItem
	err := repo.db.OrmInstance.
		Preload("Category").
		Preload("Restaurant").
		Preload("Addons").
		Where("restaurant_id = ?", restaurantID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (repo *MenuItemRepository) GetMenuItemsByCategoryID(categoryID uuid.UUID) ([]MenuItem, error) {
	var items []MenuItem
	err := repo.db.OrmInstance.
		Preload("Category").
		Preload("Restaurant").
		Preload("Addons").
		Where("category_id = ?", categoryID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (repo *MenuItemRepository) GetMenuItemByID(id uuid.UUID) (*MenuItem, error) {
	var item MenuItem
	err := repo.db.OrmInstance.
		Preload("Category").
		Preload("Restaurant").
		Preload("Addons").
		First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *MenuItemRepository) GetMenuItemsByIDs(ids []uuid.UUID) ([]MenuItem, error) {
	var items []MenuItem
	err := repo.db.OrmInstance.Preload("Addons").Where("id IN ?", ids).Find(&items).Error
	return items, err
}

func (repo *MenuItemRepository) UpdateMenuItem(item *MenuItem) error {
	return repo.db.OrmInstance.Save(item).Error
}

func (repo *MenuItemRepository) DeleteMenuItem(id uuid.UUID) error {
	return repo.db.OrmInstance.Delete(&MenuItem{}, "id = ?", id).Error
}

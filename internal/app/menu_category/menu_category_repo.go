package menucategory

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ktmbeestech/yanshi/database"
	"gorm.io/gorm"
)

type MenuCategoryRepository struct {
	db *database.OrmDb
}

func NewMenuCategoryRepository(db *database.OrmDb) *MenuCategoryRepository {
	return &MenuCategoryRepository{db: db}
}

func (repo *MenuCategoryRepository) CreateCategory(c *MenuCategory) error {
	return repo.db.OrmInstance.Create(c).Error
}

func (repo *MenuCategoryRepository) GetCategoriesByRestaurantID(restaurantID uuid.UUID) ([]MenuCategory, error) {
	var categories []MenuCategory
	err := repo.db.OrmInstance.
		Preload("Restaurant").
		Where("restaurant_id = ?", restaurantID).
		Order("created_at ASC").
		Find(&categories).Error
	return categories, err
}

func (repo *MenuCategoryRepository) GetCategoryByID(id uuid.UUID) (*MenuCategory, error) {
	var c MenuCategory
	err := repo.db.OrmInstance.
		Preload("Restaurant").
		First(&c, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (repo *MenuCategoryRepository) FindByRestaurantAndName(restaurantID uuid.UUID, name string) (*MenuCategory, error) {
	var c MenuCategory
	err := repo.db.OrmInstance.
		Where("restaurant_id = ? AND name = ?", restaurantID, name).
		First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (repo *MenuCategoryRepository) UpdateCategory(c *MenuCategory) error {
	return repo.db.OrmInstance.Save(c).Error
}

func (repo *MenuCategoryRepository) DeleteCategory(id uuid.UUID) error {
	return repo.db.OrmInstance.Delete(&MenuCategory{}, "id = ?", id).Error
}

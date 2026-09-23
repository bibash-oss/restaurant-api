package table

import (
	"errors"

	"github.com/google/uuid"
	"kitchen-api/internal/database"
	"gorm.io/gorm"
)

type TableRepository struct {
	db *database.OrmDb
}

func NewTableRepository(db *database.OrmDb) *TableRepository {
	return &TableRepository{db: db}
}

func (repo *TableRepository) CreateTable(t *Table) error {
	return repo.db.OrmInstance.Create(t).Error
}

func (repo *TableRepository) GetTablesByRestaurantID(restaurantID uuid.UUID) ([]Table, error) {
	var tables []Table
	err := repo.db.OrmInstance.
		Where("restaurant_id = ?", restaurantID).
		Order("number ASC").
		Find(&tables).Error
	return tables, err
}

func (repo *TableRepository) GetTableByID(id uuid.UUID) (*Table, error) {
	var t Table
	err := repo.db.OrmInstance.
		Preload("Restaurant").
		First(&t, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (repo *TableRepository) FindByRestaurantAndNumber(restaurantID uuid.UUID, number string) (*Table, error) {
	var t Table
	err := repo.db.OrmInstance.
		Where("restaurant_id = ? AND number = ?", restaurantID, number).
		First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (repo *TableRepository) UpdateTable(t *Table) error {
	return repo.db.OrmInstance.Save(t).Error
}

func (repo *TableRepository) DeleteTable(id uuid.UUID) error {
	return repo.db.OrmInstance.Delete(&Table{}, "id = ?", id).Error
}

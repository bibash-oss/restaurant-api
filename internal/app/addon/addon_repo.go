package addon

import (
	"errors"

	"github.com/google/uuid"
	"kitchen-api/internal/database"
	"gorm.io/gorm"
)

type AddonRepository struct {
	db *database.OrmDb
}

func NewAddonRepository(db *database.OrmDb) *AddonRepository {
	return &AddonRepository{db: db}
}

func (repo *AddonRepository) CreateAddon(a *Addon) error {
	return repo.db.OrmInstance.Create(a).Error
}

func (repo *AddonRepository) GetAddonsByRestaurantID(restaurantID uuid.UUID) ([]Addon, error) {
	var addons []Addon
	err := repo.db.OrmInstance.
		Where("restaurant_id = ?", restaurantID).
		Order("name ASC").
		Find(&addons).Error
	return addons, err
}

func (repo *AddonRepository) GetAddonByID(id uuid.UUID) (*Addon, error) {
	var a Addon
	err := repo.db.OrmInstance.
		Preload("Restaurant").
		First(&a, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (repo *AddonRepository) GetAddonsByIDs(ids []uuid.UUID) ([]Addon, error) {
	var addons []Addon
	err := repo.db.OrmInstance.Where("id IN ?", ids).Find(&addons).Error
	return addons, err
}

func (repo *AddonRepository) FindByRestaurantAndName(restaurantID uuid.UUID, name string) (*Addon, error) {
	var a Addon
	err := repo.db.OrmInstance.
		Where("restaurant_id = ? AND name = ?", restaurantID, name).
		First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (repo *AddonRepository) UpdateAddon(a *Addon) error {
	return repo.db.OrmInstance.Save(a).Error
}

func (repo *AddonRepository) DeleteAddon(id uuid.UUID) error {
	return repo.db.OrmInstance.Delete(&Addon{}, "id = ?", id).Error
}

type menuItemAddonRow struct {
	MenuItemID uuid.UUID `gorm:"column:menu_item_id;primaryKey"`
	AddonID    uuid.UUID `gorm:"column:addon_id;primaryKey"`
}

func (menuItemAddonRow) TableName() string {
	return "menu_item_addons"
}

func (repo *AddonRepository) GetMenuItemRestaurantID(menuItemID uuid.UUID) (*uuid.UUID, error) {
	var row struct {
		RestaurantID uuid.UUID `gorm:"column:restaurant_id"`
	}
	err := repo.db.OrmInstance.Table("menu_items").Select("restaurant_id").Where("id = ?", menuItemID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row.RestaurantID, nil
}

func (repo *AddonRepository) AssignAddonsToMenuItem(menuItemID uuid.UUID, addonIDs []uuid.UUID) error {
	return repo.db.OrmInstance.Transaction(func(tx *gorm.DB) error {
		// Remove existing mappings for this menu item
		if err := tx.Where("menu_item_id = ?", menuItemID).Delete(&menuItemAddonRow{}).Error; err != nil {
			return err
		}

		// Insert new mappings
		if len(addonIDs) > 0 {
			mappings := make([]menuItemAddonRow, len(addonIDs))
			for i, aid := range addonIDs {
				mappings[i] = menuItemAddonRow{
					MenuItemID: menuItemID,
					AddonID:    aid,
				}
			}
			if err := tx.Create(&mappings).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *AddonRepository) GetAddonsByMenuItemID(menuItemID uuid.UUID) ([]Addon, error) {
	var addons []Addon
	err := repo.db.OrmInstance.
		Table("addons").
		Joins("JOIN menu_item_addons ON menu_item_addons.addon_id = addons.id").
		Where("menu_item_addons.menu_item_id = ? AND addons.is_active = ?", menuItemID, true).
		Order("addons.name ASC").
		Find(&addons).Error
	return addons, err
}

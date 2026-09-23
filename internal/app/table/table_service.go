package table

import (
	"fmt"

	"github.com/google/uuid"
)

type TableService struct {
	tableRepo *TableRepository
}

func NewTableService(tableRepo *TableRepository) *TableService {
	return &TableService{tableRepo: tableRepo}
}

func (s *TableService) CreateTable(req *CreateTableRequest) (*Table, error) {
	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	existing, err := s.tableRepo.FindByRestaurantAndNumber(restaurantID, req.Number)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("table number '%s' already exists for this restaurant", req.Number)
	}

	t := &Table{
		RestaurantID: restaurantID,
		Number:       req.Number,
		IsActive:     true,
	}

	if err := s.tableRepo.CreateTable(t); err != nil {
		return nil, err
	}

	createdTable, err := s.tableRepo.GetTableByID(t.ID)
	if err == nil && createdTable != nil {
		return createdTable, nil
	}

	return t, nil
}

func (s *TableService) GetTablesByRestaurantID(restaurantIDStr string) ([]Table, error) {
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}
	return s.tableRepo.GetTablesByRestaurantID(restaurantID)
}

func (s *TableService) GetTableByID(idStr string) (*Table, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid table ID")
	}
	t, err := s.tableRepo.GetTableByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("table not found")
	}
	return t, nil
}

func (s *TableService) UpdateTable(idStr string, req *UpdateTableRequest) (*Table, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid table ID")
	}
	t, err := s.tableRepo.GetTableByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("table not found")
	}

	if req.Number != nil {
		// check uniqueness if changing number
		if *req.Number != t.Number {
			existing, err := s.tableRepo.FindByRestaurantAndNumber(t.RestaurantID, *req.Number)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				return nil, fmt.Errorf("table number '%s' already exists for this restaurant", *req.Number)
			}
			t.Number = *req.Number
		}
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}

	if err := s.tableRepo.UpdateTable(t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TableService) DeleteTable(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid table ID")
	}
	return s.tableRepo.DeleteTable(id)
}

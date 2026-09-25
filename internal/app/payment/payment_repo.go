package payment

import (
	"errors"

	"gorm.io/gorm"
	"kitchen-api/internal/database"
)

type PaymentRepository struct {
	db *database.OrmDb
}

func NewPaymentRepository(db *database.OrmDb) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreateSession(session *PaymentSession) error {
	return r.db.OrmInstance.Create(session).Error
}

func (r *PaymentRepository) GetByStripeSessionID(sessionID string) (*PaymentSession, error) {
	var session PaymentSession
	err := r.db.OrmInstance.
		Preload("Order").
		Preload("Order.OrderItems").
		Preload("Order.Table").
		Preload("Order.Restaurant").
		Where("stripe_session_id = ?", sessionID).
		First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *PaymentRepository) UpdateSession(session *PaymentSession) error {
	return r.db.OrmInstance.Save(session).Error
}

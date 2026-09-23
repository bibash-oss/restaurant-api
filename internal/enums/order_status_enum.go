package enums

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusReady     OrderStatus = "READY"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending,
		OrderStatusConfirmed,
		OrderStatusPreparing,
		OrderStatusReady,
		OrderStatusCompleted,
		OrderStatusCancelled:
		return true
	default:
		return false
	}
}

package order

import (
	"fmt"
	"math"
	"strings"

	"kitchen-api/internal/app/addon"
	menuitem "kitchen-api/internal/app/menu_item"
	orderitem "kitchen-api/internal/app/order_item"
	"kitchen-api/internal/app/table"
	"kitchen-api/internal/enums"

	"github.com/google/uuid"
)

type LineItemDetail struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
}

type ValidatedOrderData struct {
	RestaurantID uuid.UUID
	TableID      uuid.UUID
	TotalAmount  float64
	LineItems    []LineItemDetail
	OrderItems   []orderitem.OrderItem
}

type OrderService struct {
	orderRepo *OrderRepository
	tableRepo *table.TableRepository
	itemRepo  *menuitem.MenuItemRepository
	addonRepo *addon.AddonRepository
}

func NewOrderService(
	orderRepo *OrderRepository,
	tableRepo *table.TableRepository,
	itemRepo *menuitem.MenuItemRepository,
	addonRepo *addon.AddonRepository,
) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		tableRepo: tableRepo,
		itemRepo:  itemRepo,
		addonRepo: addonRepo,
	}
}

func (s *OrderService) buildAndValidateOrder(req *CreateOrderRequest) (*ValidatedOrderData, error) {
	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	tableID, err := uuid.Parse(req.TableID)
	if err != nil {
		return nil, fmt.Errorf("invalid table ID")
	}

	// Validate table exists, is active, and belongs to this restaurant
	tbl, err := s.tableRepo.GetTableByID(tableID)
	if err != nil {
		return nil, err
	}
	if tbl == nil {
		return nil, fmt.Errorf("table not found")
	}
	if tbl.RestaurantID != restaurantID {
		return nil, fmt.Errorf("table does not belong to this restaurant")
	}
	if !tbl.IsActive {
		return nil, fmt.Errorf("table is currently inactive")
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	// 1. Collect unique menu item IDs and addon IDs
	uniqueMenuItemIDs := make(map[uuid.UUID]bool)
	uniqueAddonIDs := make(map[uuid.UUID]bool)

	type parsedAddon struct {
		addonID  uuid.UUID
		quantity int
	}

	type parsedItem struct {
		menuItemID uuid.UUID
		quantity   int
		addons     []parsedAddon
	}
	parsedItems := make([]parsedItem, len(req.Items))

	for i, itemReq := range req.Items {
		mID, err := uuid.Parse(itemReq.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("invalid menu item ID: %s", itemReq.MenuItemID)
		}
		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0")
		}
		uniqueMenuItemIDs[mID] = true

		pAddons := make([]parsedAddon, 0, len(itemReq.Addons))
		for _, aReq := range itemReq.Addons {
			aID, err := uuid.Parse(aReq.AddonID)
			if err != nil {
				return nil, fmt.Errorf("invalid addon ID: %s", aReq.AddonID)
			}
			if aReq.Quantity <= 0 {
				return nil, fmt.Errorf("addon quantity must be greater than 0")
			}
			uniqueAddonIDs[aID] = true
			pAddons = append(pAddons, parsedAddon{
				addonID:  aID,
				quantity: aReq.Quantity,
			})
		}

		parsedItems[i] = parsedItem{
			menuItemID: mID,
			quantity:   itemReq.Quantity,
			addons:     pAddons,
		}
	}

	// 2. Fetch and validate menu items
	menuItemIDsSlice := make([]uuid.UUID, 0, len(uniqueMenuItemIDs))
	for id := range uniqueMenuItemIDs {
		menuItemIDsSlice = append(menuItemIDsSlice, id)
	}

	menuItems, err := s.itemRepo.GetMenuItemsByIDs(menuItemIDsSlice)
	if err != nil {
		return nil, err
	}
	if len(menuItems) != len(uniqueMenuItemIDs) {
		return nil, fmt.Errorf("one or more menu items were not found")
	}

	menuItemMap := make(map[uuid.UUID]menuitem.MenuItem, len(menuItems))
	for _, mi := range menuItems {
		if mi.RestaurantID != restaurantID {
			return nil, fmt.Errorf("menu item '%s' does not belong to this restaurant", mi.Name)
		}
		if !mi.IsActive {
			return nil, fmt.Errorf("menu item '%s' is currently unavailable", mi.Name)
		}
		menuItemMap[mi.ID] = mi
	}

	// 3. Fetch and validate addons (if any requested)
	addonMap := make(map[uuid.UUID]addon.Addon)
	if len(uniqueAddonIDs) > 0 {
		addonIDsSlice := make([]uuid.UUID, 0, len(uniqueAddonIDs))
		for id := range uniqueAddonIDs {
			addonIDsSlice = append(addonIDsSlice, id)
		}

		addons, err := s.addonRepo.GetAddonsByIDs(addonIDsSlice)
		if err != nil {
			return nil, err
		}
		if len(addons) != len(uniqueAddonIDs) {
			return nil, fmt.Errorf("one or more addons were not found")
		}

		for _, a := range addons {
			if a.RestaurantID != restaurantID {
				return nil, fmt.Errorf("addon '%s' does not belong to this restaurant", a.Name)
			}
			if !a.IsActive {
				return nil, fmt.Errorf("addon '%s' is currently unavailable", a.Name)
			}
			addonMap[a.ID] = a
		}
	}

	// 4. Build order items, validate addons per menu item, and calculate total
	var totalAmount float64
	orderItems := make([]orderitem.OrderItem, 0, len(parsedItems))
	lineItems := make([]LineItemDetail, 0, len(parsedItems))

	for _, pi := range parsedItems {
		mi := menuItemMap[pi.menuItemID]

		// Build set of allowed addons for this menu item
		allowedAddonIDs := make(map[uuid.UUID]bool, len(mi.Addons))
		for _, a := range mi.Addons {
			allowedAddonIDs[a.ID] = true
		}

		var itemAddonTotal float64
		itemAddons := make([]orderitem.OrderItemAddon, 0, len(pi.addons))
		addonDescriptions := make([]string, 0, len(pi.addons))

		orderItemID := uuid.New()
		for _, pa := range pi.addons {
			a, ok := addonMap[pa.addonID]
			if !ok {
				return nil, fmt.Errorf("addon not found")
			}
			if !allowedAddonIDs[pa.addonID] {
				return nil, fmt.Errorf("addon '%s' is not available for menu item '%s'", a.Name, mi.Name)
			}

			itemAddonTotal += a.Price * float64(pa.quantity)
			addonDescriptions = append(addonDescriptions, fmt.Sprintf("%s (x%d)", a.Name, pa.quantity))
			itemAddons = append(itemAddons, orderitem.OrderItemAddon{
				ID:          uuid.New(),
				OrderItemID: orderItemID,
				AddonID:     a.ID,
				Quantity:    pa.quantity,
				UnitPrice:   a.Price,
			})
		}

		lineItemUnitPrice := mi.Price + itemAddonTotal
		totalAmount += lineItemUnitPrice * float64(pi.quantity)

		orderItems = append(orderItems, orderitem.OrderItem{
			ID:         orderItemID,
			MenuItemID: mi.ID,
			Quantity:   pi.quantity,
			UnitPrice:  mi.Price,
			Addons:     itemAddons,
		})

		lineItems = append(lineItems, LineItemDetail{
			Name:        mi.Name,
			Description: strings.Join(addonDescriptions, ", "),
			Quantity:    pi.quantity,
			UnitPrice:   lineItemUnitPrice,
		})
	}

	// Round totalAmount to 2 decimal places
	totalAmount = math.Round(totalAmount*100) / 100

	return &ValidatedOrderData{
		RestaurantID: restaurantID,
		TableID:      tableID,
		TotalAmount:  totalAmount,
		LineItems:    lineItems,
		OrderItems:   orderItems,
	}, nil
}

func (s *OrderService) ValidateOrderForPayment(req *CreateOrderRequest) (*ValidatedOrderData, error) {
	return s.buildAndValidateOrder(req)
}

func (s *OrderService) CreatePaidOrder(req *CreateOrderRequest, stripeSessionID string) (*Order, error) {
	// Idempotency: Return existing order if already created for this Stripe session
	if stripeSessionID != "" {
		existing, err := s.orderRepo.GetOrderByStripeSessionID(stripeSessionID)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	validated, err := s.buildAndValidateOrder(req)
	if err != nil {
		return nil, err
	}

	newOrder := &Order{
		RestaurantID:    validated.RestaurantID,
		TableID:         validated.TableID,
		Status:          enums.OrderStatusPending,
		TotalAmount:     validated.TotalAmount,
		Notes:           req.Notes,
		PaymentStatus:   "PAID",
		StripeSessionID: &stripeSessionID,
	}

	if err := s.orderRepo.CreateOrderWithItems(newOrder, validated.OrderItems); err != nil {
		return nil, err
	}

	fullOrder, err := s.orderRepo.GetOrderByID(newOrder.ID)
	if err == nil && fullOrder != nil {
		return fullOrder, nil
	}

	return newOrder, nil
}

func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*Order, error) {
	validated, err := s.buildAndValidateOrder(req)
	if err != nil {
		return nil, err
	}

	newOrder := &Order{
		RestaurantID:  validated.RestaurantID,
		TableID:       validated.TableID,
		Status:        enums.OrderStatusPending,
		TotalAmount:   validated.TotalAmount,
		Notes:         req.Notes,
		PaymentStatus: "UNPAID",
	}

	if err := s.orderRepo.CreateOrderWithItems(newOrder, validated.OrderItems); err != nil {
		return nil, err
	}

	// Return preloaded order with Table, Restaurant, and Item details for receipt printing
	fullOrder, err := s.orderRepo.GetOrderByID(newOrder.ID)
	if err == nil && fullOrder != nil {
		return fullOrder, nil
	}

	return newOrder, nil
}

func (s *OrderService) GetOrdersByRestaurantID(restaurantIDStr string) ([]Order, error) {
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}
	return s.orderRepo.GetOrdersByRestaurantID(restaurantID)
}

func (s *OrderService) GetOrdersByTableID(tableIDStr string) ([]Order, error) {
	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid table ID")
	}
	return s.orderRepo.GetOrdersByTableID(tableID)
}

func (s *OrderService) GetOrderByID(idStr string) (*Order, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID")
	}
	o, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, fmt.Errorf("order not found")
	}
	return o, nil
}

func (s *OrderService) UpdateOrderStatus(idStr string, req *UpdateOrderStatusRequest) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid order ID")
	}

	if !req.Status.IsValid() {
		return fmt.Errorf("invalid order status '%s'", req.Status)
	}

	o, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return err
	}
	if o == nil {
		return fmt.Errorf("order not found")
	}

	return s.orderRepo.UpdateOrderStatus(id, req.Status)
}

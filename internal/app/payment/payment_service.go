package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/webhook"
	"kitchen-api/internal/app/order"
	"kitchen-api/internal/config"
)

type PaymentService struct {
	paymentRepo  *PaymentRepository
	orderService *order.OrderService
	mu           sync.Mutex
}

func NewPaymentService(paymentRepo *PaymentRepository, orderService *order.OrderService) *PaymentService {
	return &PaymentService{
		paymentRepo:  paymentRepo,
		orderService: orderService,
	}
}

func (s *PaymentService) getStripeKey() string {
	key := os.Getenv("STRIPE_SECRET_KEY")
	fmt.Println("stripe key", key)
	if key != "" {
		return key
	}

	return config.Default().GetString("stripe.secretKey")
}

func (s *PaymentService) getWebhookSecret() string {
	if secret := os.Getenv("STRIPE_WEBHOOK_SECRET"); secret != "" {
		return secret
	}
	return config.Default().GetString("stripe.webhookSecret")
}

func (s *PaymentService) getCurrency() string {
	if cur := os.Getenv("STRIPE_CURRENCY"); cur != "" {
		return strings.ToLower(cur)
	}
	cur := config.Default().GetString("stripe.currency")
	if cur == "" {
		cur = "aud"
	}
	return strings.ToLower(cur)
}

func (s *PaymentService) getAppURL() string {
	appURL := config.Default().GetString("app.url")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	return appURL
}

func (s *PaymentService) CreateCheckoutSession(req *order.CreateOrderRequest) (*CreateCheckoutResponse, error) {
	stripeKey := s.getStripeKey()
	if stripeKey == "" {
		return nil, fmt.Errorf("stripe secret key is not configured")
	}
	stripe.Key = stripeKey

	// 1. Validate order and calculate total securely server-side
	validated, err := s.orderService.ValidateOrderForPayment(req)
	if err != nil {
		return nil, err
	}

	currency := s.getCurrency()
	appURL := s.getAppURL()

	// 2. Build Stripe line items
	lineItemParams := make([]*stripe.CheckoutSessionLineItemParams, 0, len(validated.LineItems))
	for _, item := range validated.LineItems {
		unitAmountCents := int64(math.Round(item.UnitPrice * 100))
		productData := &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
			Name: stripe.String(item.Name),
		}
		if item.Description != "" {
			productData.Description = stripe.String(item.Description)
		}

		lineItemParams = append(lineItemParams, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:    stripe.String(currency),
				ProductData: productData,
				UnitAmount:  stripe.Int64(unitAmountCents),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	// 3. Serialize payload so order can be created on successful payment
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize order payload: %w", err)
	}

	// 4. Create Stripe Checkout Session
	successURL := fmt.Sprintf("%s/order/success?session_id={CHECKOUT_SESSION_ID}", appURL)
	cancelURL := fmt.Sprintf("%s/order/cancel", appURL)

	sessionParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems:  lineItemParams,
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Locale:     stripe.String("auto"),
		Metadata: map[string]string{
			"restaurant_id": req.RestaurantID,
			"table_id":      req.TableID,
		},
	}

	// Pre-fill customer address with Australia so the Country dropdown defaults to Australia
	cust, err := customer.New(&stripe.CustomerParams{
		Address: &stripe.AddressParams{
			Country: stripe.String("AU"),
		},
	})
	if err == nil && cust != nil {
		sessionParams.Customer = stripe.String(cust.ID)
		sessionParams.CustomerUpdate = &stripe.CheckoutSessionCustomerUpdateParams{
			Address: stripe.String("auto"),
		}
	}

	sess, err := session.New(sessionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe checkout session: %w", err)
	}

	// 5. Store pending payment session in DB
	paymentSession := &PaymentSession{
		StripeSessionID: sess.ID,
		RestaurantID:    validated.RestaurantID,
		TableID:         validated.TableID,
		OrderPayload:    string(payloadBytes),
		TotalAmount:     validated.TotalAmount,
		Currency:        currency,
		Status:          "PENDING",
	}

	if err := s.paymentRepo.CreateSession(paymentSession); err != nil {
		return nil, fmt.Errorf("failed to save payment session: %w", err)
	}

	return &CreateCheckoutResponse{
		CheckoutURL: sess.URL,
		SessionID:   sess.ID,
		TotalAmount: validated.TotalAmount,
	}, nil
}

// completePaymentSession idempotently creates the order and marks the payment session COMPLETED
func (s *PaymentService) completePaymentSession(sessionID string) (*order.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionRecord, err := s.paymentRepo.GetByStripeSessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("database error fetching payment session: %w", err)
	}
	if sessionRecord == nil {
		return nil, fmt.Errorf("payment session not found for id %s", sessionID)
	}

	// Idempotency: if order was already created, return it immediately
	if sessionRecord.Status == "COMPLETED" && sessionRecord.OrderID != nil {
		if sessionRecord.Order != nil {
			return sessionRecord.Order, nil
		}
		return s.orderService.GetOrderByID(sessionRecord.OrderID.String())
	}

	// Unmarshal the saved order request payload
	var orderReq order.CreateOrderRequest
	if err := json.Unmarshal([]byte(sessionRecord.OrderPayload), &orderReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order payload: %w", err)
	}

	// Create the order in orders table (idempotent at repository level as well)
	createdOrder, err := s.orderService.CreatePaidOrder(&orderReq, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order after payment: %w", err)
	}

	// Update payment session to COMPLETED
	sessionRecord.Status = "COMPLETED"
	sessionRecord.OrderID = &createdOrder.ID
	sessionRecord.Order = createdOrder
	if err := s.paymentRepo.UpdateSession(sessionRecord); err != nil {
		return nil, fmt.Errorf("failed to update payment session: %w", err)
	}

	return createdOrder, nil
}

func (s *PaymentService) HandleWebhook(payload []byte, sigHeader string) error {
	var event stripe.Event
	webhookSecret := s.getWebhookSecret()

	if webhookSecret != "" {
		var err error
		event, err = webhook.ConstructEventWithOptions(payload, sigHeader, webhookSecret, webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
		if err != nil {
			return fmt.Errorf("invalid webhook signature: %w", err)
		}
	} else {
		// If secret is not configured, parse payload directly
		if err := json.Unmarshal(payload, &event); err != nil {
			return fmt.Errorf("failed to parse webhook event: %w", err)
		}
	}

	if event.Type == "checkout.session.completed" {
		var sessionID string
		var checkoutSess stripe.CheckoutSession
		if len(event.Data.Raw) > 0 {
			if err := json.Unmarshal(event.Data.Raw, &checkoutSess); err == nil {
				sessionID = checkoutSess.ID
			}
		}
		if sessionID == "" && event.Data.Object != nil {
			if idVal, ok := event.Data.Object["id"].(string); ok {
				sessionID = idVal
			}
		}
		if sessionID == "" {
			return fmt.Errorf("failed to extract checkout session id from event")
		}

		if _, err := s.completePaymentSession(sessionID); err != nil {
			return err
		}
	}

	return nil
}

func (s *PaymentService) GetSessionStatus(sessionID string) (*SessionStatusResponse, error) {
	sessionRecord, err := s.paymentRepo.GetByStripeSessionID(sessionID)
	if err != nil {
		return nil, err
	}
	if sessionRecord == nil {
		return nil, errors.New("session not found")
	}

	// Fallback sync: If webhook hasn't processed it yet, verify directly with Stripe API
	if sessionRecord.Status != "COMPLETED" || sessionRecord.OrderID == nil {
		stripeKey := s.getStripeKey()
		if stripeKey != "" {
			stripe.Key = stripeKey
			stripeSess, err := session.Get(sessionID, nil)
			if err == nil && stripeSess != nil {
				if stripeSess.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid || stripeSess.Status == stripe.CheckoutSessionStatusComplete {
					if createdOrder, err := s.completePaymentSession(sessionID); err == nil {
						sessionRecord.Status = "COMPLETED"
						sessionRecord.OrderID = &createdOrder.ID
						sessionRecord.Order = createdOrder
					}
				}
			}
		}
	}

	var orderIDStr *string
	if sessionRecord.OrderID != nil {
		str := sessionRecord.OrderID.String()
		orderIDStr = &str
	}

	return &SessionStatusResponse{
		SessionID: sessionRecord.StripeSessionID,
		Status:    sessionRecord.Status,
		OrderID:   orderIDStr,
		Order:     sessionRecord.Order,
	}, nil
}

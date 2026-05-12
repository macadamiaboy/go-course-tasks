package main

import (
	"context"
	"errors"
	"fmt"
)

type CreateOrderRequest struct {
	CustomerID string
	Items      []string
}

type CreateOrderResponse struct {
	OrderID string
	Status  string
}

type GetOrderRequest struct {
	OrderID string
}

type GetOrderResponse struct {
	OrderID    string
	CustomerID string
	Items      []string
	Status     string
}

var ErrOrderNotFound = errors.New("order not found")
var ErrNoDataProvided = errors.New("there's no provided data")

type OrderServiceServer interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error)
}

// TODO: Implement orderServiceImpl.
// - Store orders in a map keyed by order ID.
// - Generate IDs like "order-1", "order-2", etc.
// - CreateOrder: validate CustomerID and Items are non-empty; store and return status "created".
// - GetOrder: return ErrOrderNotFound for unknown IDs.

type order struct {
	id         string
	customerID string
	items      []string
	status     string
}

type orderServiceImpl struct {
	orders map[string]order
	nextID int
}

func NewOrderService() OrderServiceServer {
	return &orderServiceImpl{
		orders: make(map[string]order),
		nextID: 1,
	}
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// TODO: validate input, store order, return response
	if req.CustomerID == "" {
		return &CreateOrderResponse{Status: ErrNoDataProvided.Error()}, fmt.Errorf("customer id: %w", ErrNoDataProvided)
	}

	if req.Items == nil {
		return &CreateOrderResponse{Status: ErrNoDataProvided.Error()}, fmt.Errorf("items: %w", ErrNoDataProvided)
	}

	orderId := fmt.Sprintf("order-%d", s.nextID)
	order := order{id: orderId, customerID: req.CustomerID, items: req.Items, status: "created"}

	s.orders[orderId] = order
	s.nextID++

	response := CreateOrderResponse{Status: "created", OrderID: orderId}

	return &response, nil
}

func (s *orderServiceImpl) GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	// TODO: look up order by ID, return ErrOrderNotFound if missing
	order, ok := s.orders[req.OrderID]
	if !ok {
		return &GetOrderResponse{Status: ErrOrderNotFound.Error()}, ErrOrderNotFound
	}

	return &GetOrderResponse{OrderID: order.id, CustomerID: order.customerID, Items: order.items, Status: "ok"}, nil
}

func main() {
	svc := NewOrderService()
	ctx := context.Background()

	created, err := svc.CreateOrder(ctx, &CreateOrderRequest{
		CustomerID: "c-1",
		Items:      []string{"item-a", "item-b"},
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Created: %+v\n", created)

	got, err := svc.GetOrder(ctx, &GetOrderRequest{OrderID: created.OrderID})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Got: %+v\n", got)
	}

	_, err = svc.GetOrder(ctx, &GetOrderRequest{OrderID: "unknown"})
	if errors.Is(err, ErrOrderNotFound) {
		fmt.Println("Order 'unknown': not found (expected)")
	}
}

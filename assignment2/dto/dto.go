package dto

import "time"

type CreateOrderRequest struct {
	CustomerName string                   `json:"customer_name"`
	Items        []CreateOrderItemRequest `json:"items"`
}

type CreateOrderResponse struct {
	ID           string       `json:"order_id"`
	CustomerName string       `json:"customer_name"`
	OrderAt      time.Time    `json:"order_at"`
	Items        []CreateItem `json:"items"`
}

type UpdateOrderRequest struct {
	CustomerName string                   `json:"customer_name"`
	Items        []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type CreateItem struct {
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type CreateItemResponse struct {
	ID          string `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

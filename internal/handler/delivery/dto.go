package delivery

import "time"

type DeliveryAssign struct {
	CourierID     int64     `json:"courier_id"`
	OrderID       string    `json:"order_id"`
	TransportType string    `json:"transport_type"`
	Deadline      time.Time `json:"delivery_deadline"`
}

type OrderID struct {
	OrderID string `json:"order_id"`
}

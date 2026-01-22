package delivery

import "time"

type AssignmentDelivery struct {
	CourierID int64
	OrderID   string
	TransportType string
	Deadline  time.Time
}

type Courier struct {
	ID            int64
	TransportType string
}

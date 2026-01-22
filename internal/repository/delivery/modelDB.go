package delivery

import "time"

type DeliveryDB struct {
	ID         int64
	CourierID int64
	OrderID   string
	AssignedAt time.Time
	Deadline   time.Time
}

type CourierDB struct {
	ID            int64
	Name          string
	Phone         string
	Status        string
	TransportType string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

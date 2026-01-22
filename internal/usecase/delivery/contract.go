package delivery

import (
	"context"

	"github.com/totorialman/test_project_go/internal/repository/courier"
	"github.com/totorialman/test_project_go/internal/repository/delivery"
)

type courierRepository interface {
	GetAvailable(ctx context.Context) (courier.CourierDB, error)
}

type deliveryRepository interface {
	Create(ctx context.Context, d delivery.DeliveryDB) error
}

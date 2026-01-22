package delivery

import (
	"context"
	"time"

	"github.com/totorialman/test_project_go/internal/repository/delivery"
)

type Usecase struct {
	courierRepo  courierRepository
	deliveryRepo deliveryRepository
}

func NewUsecase(courierRepo courierRepository, deliveryRepo deliveryRepository) *Usecase {
	return &Usecase{
		courierRepo:  courierRepo,
		deliveryRepo: deliveryRepo,
	}
}

func (s *Usecase) Create(ctx context.Context, orderID string) (AssignmentDelivery, error) {
	dbCourier, err := s.courierRepo.GetAvailable(ctx)
	if err != nil {
		return AssignmentDelivery{}, err
	}

	dbDelivery := delivery.DeliveryDB{
		CourierID: dbCourier.ID,
		OrderID:   orderID,
		Deadline:  time.Now().Add(30 * time.Second),
	}

	err = s.deliveryRepo.Create(ctx, dbDelivery)
	if err != nil {
		return AssignmentDelivery{}, err
	}

	delivery := AssignmentDelivery{
		CourierID:     dbDelivery.CourierID,
		OrderID:       dbDelivery.OrderID,
		TransportType: dbCourier.TransportType,
		Deadline:      dbDelivery.Deadline,
	}

	return delivery, nil
}

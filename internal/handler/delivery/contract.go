package delivery

import (
	"context"

	"github.com/totorialman/test_project_go/internal/usecase/delivery"
)

type usecase interface {
	Create(ctx context.Context, orderID string) (delivery.AssignmentDelivery, error)
}

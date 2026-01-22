package courier

import (
	"context"

	"github.com/totorialman/test_project_go/internal/repository/courier"
)

type repository interface {
	Create(ctx context.Context, c courier.CourierDB) error
	GetByPhone(ctx context.Context, phone string) (courier.CourierDB, error)
	GetAll(ctx context.Context) ([]courier.CourierDB, error)
	GetById(ctx context.Context, id int64) (courier.CourierDB, error)
	Update(ctx context.Context, c courier.CourierDB) error
}

package courier

import (
	"context"

	"github.com/totorialman/test_project_go/internal/usecase/courier"
)

type usecase interface {
	Create(ctx context.Context, c courier.Courier) error
	GetAll(ctx context.Context) ([]courier.Courier, error)
	GetById(ctx context.Context, id int64) (courier.Courier, error)
	Update(ctx context.Context, c courier.Courier) error
}

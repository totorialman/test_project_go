package delivery

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, d DeliveryDB) error {
	const createDeliveryQuery = `
	INSERT INTO delivery(courier_id, order_id, deadline)
	VALUES($1, $2, $3);
	`
	_, err := r.db.Exec(ctx, createDeliveryQuery, d.CourierID, d.OrderID, d.Deadline)
	if err != nil {
		return fmt.Errorf("falied to create delivery: %w", err)
	}

	return nil
}

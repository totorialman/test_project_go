package courier

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	courierErrors "github.com/totorialman/test_project_go/internal/errors/courier"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]CourierDB, error) {
	const getCouriersQuery = `
	SELECT id, name, phone, status, transport_type, created_at, updated_at FROM couriers;
	`

	rows, err := r.db.Query(ctx, getCouriersQuery)
	if err != nil {
		return []CourierDB{}, fmt.Errorf("failed to query couriers: %w", err)
	}
	defer rows.Close()

	var dbCouriers []CourierDB
	for rows.Next() {
		var dbCourier CourierDB
		err := rows.Scan(
			&dbCourier.ID,
			&dbCourier.Name,
			&dbCourier.Phone,
			&dbCourier.Status,
			&dbCourier.TransportType,
			&dbCourier.CreatedAt,
			&dbCourier.UpdatedAt,
		)
		if err != nil {
			return []CourierDB{}, fmt.Errorf("failed to scan courier row: %w", err)
		}
		dbCouriers = append(dbCouriers, dbCourier)
	}
	if err := rows.Err(); err != nil {
		return []CourierDB{}, fmt.Errorf("error iterating courier rows: %w", err)
	}

	return dbCouriers, nil
}

func (r *Repository) Create(ctx context.Context, c CourierDB) error {
	const createCourierQuery = `
	INSERT INTO couriers(name, phone, status, transport_type) 
	VALUES($1, $2, $3, $4);
	`

	if c.TransportType == "" {
		c.TransportType = "on_foot"
	}

	_, err := r.db.Exec(ctx, createCourierQuery, c.Name, c.Phone, c.Status, c.TransportType)
	if err != nil {
		fmt.Println("failed to create courier:%w", err)
		return fmt.Errorf("failed to create courier:%w", err)
	}

	return nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (CourierDB, error) {
	var dbCourier CourierDB

	const getCourierByPhoneQuery = `
	SELECT id, name, phone, status, transport_type, created_at, updated_at FROM couriers WHERE phone = $1;
	`

	err := r.db.QueryRow(ctx, getCourierByPhoneQuery, phone).Scan(
		&dbCourier.ID,
		&dbCourier.Name,
		&dbCourier.Phone,
		&dbCourier.Status,
		&dbCourier.TransportType,
		&dbCourier.CreatedAt,
		&dbCourier.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("courier not found")
		return CourierDB{}, courierErrors.ErrCourierNotFound
	}
	if err != nil {
		fmt.Println("failed to get courier by phone: %w", err)
		return CourierDB{}, fmt.Errorf("failed to get courier by phone: %w", err)
	}

	return dbCourier, nil

}

func (r *Repository) GetById(ctx context.Context, id int64) (CourierDB, error) {
	const getCourierById = `
	SELECT id, name, phone, status, transport_type, created_at, updated_at FROM couriers WHERE id=$1;
	`
	var dbCourier CourierDB
	err := r.db.QueryRow(ctx, getCourierById, id).Scan(
		&dbCourier.ID,
		&dbCourier.Name,
		&dbCourier.Phone,
		&dbCourier.Status,
		&dbCourier.TransportType,
		&dbCourier.CreatedAt,
		&dbCourier.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return CourierDB{}, courierErrors.ErrCourierNotFound
	}
	if err != nil {
		return CourierDB{}, fmt.Errorf("failed to get courier by id: %w", err)
	}

	return dbCourier, nil
}

func (r *Repository) Update(ctx context.Context, c CourierDB) error {
	const updateCourierQuery = `
	UPDATE couriers
	SET
		name = COALESCE(NULLIF($1, ''), name),
		phone = COALESCE(NULLIF($2, ''), phone),
		status = COALESCE(NULLIF($3, ''), status),
		transport_type = COALESCE(NULLIF($4, ''), transport_type),
		updated_at = NOW()
	WHERE id = $5
	RETURNING id;
	`

	var id int64
	err := r.db.QueryRow(ctx, updateCourierQuery, c.Name, c.Phone, c.Status, c.TransportType, c.ID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return courierErrors.ErrCourierNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update courier: %w", err)
	}

	return nil
}

func (r *Repository) GetAvailable(ctx context.Context) (CourierDB, error) {
	const getAvailableQuery = `
	SELECT id, transport_type FROM couriers WHERE status = 'available' LIMIT 1;
	`
	var dbCourier CourierDB
	err := r.db.QueryRow(ctx, getAvailableQuery).Scan(&dbCourier.ID, &dbCourier.TransportType)
	if errors.Is(err, pgx.ErrNoRows) {
		return CourierDB{}, courierErrors.ErrCourierNotFound
	}
	if err != nil {
		return CourierDB{}, fmt.Errorf("failed to get available courier: %w", err)
	}

	return dbCourier, nil
}

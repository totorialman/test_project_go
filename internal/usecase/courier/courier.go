package courier

import (
	"context"
	"errors"
	"fmt"

	courierErrors "github.com/totorialman/test_project_go/internal/errors/courier"
	"github.com/totorialman/test_project_go/internal/repository/courier"
)

type Usecase struct {
	repo repository
}

func NewService(repo repository) *Usecase {
	return &Usecase{repo: repo}
}

func (s *Usecase) Create(ctx context.Context, c Courier) error {
	_, err := s.repo.GetByPhone(ctx, c.Phone)
	if err != nil && !errors.Is(err, courierErrors.ErrCourierNotFound) {
		fmt.Println("failed to check existing courier: %w", err)
		return fmt.Errorf("failed to check existing courier: %w", err)
	}
	if err == nil {
		fmt.Println("courier with this phone already exist")
		return courierErrors.ErrCourierExists

	}
	dbCourier := courier.CourierDB{
		Name:          c.Name,
		Phone:         c.Phone,
		Status:        c.Status,
		TransportType: c.TransportType,
	}

	return s.repo.Create(ctx, dbCourier)
}

func (s *Usecase) GetAll(ctx context.Context) ([]Courier, error) {
	dbCouriers, err := s.repo.GetAll(ctx)
	if err != nil {
		return []Courier{}, err
	}

	couriers := make([]Courier, len(dbCouriers))
	for i, c := range dbCouriers {
		couriers[i] = Courier{
			ID:            c.ID,
			Name:          c.Name,
			Phone:         c.Phone,
			Status:        c.Status,
			TransportType: c.TransportType,
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
		}
	}

	return couriers, nil
}

func (s *Usecase) GetById(ctx context.Context, id int64) (Courier, error) {
	dbCourier, err := s.repo.GetById(ctx, id)
	if errors.Is(err, courierErrors.ErrCourierNotFound) {
		return Courier{}, err
	}
	if err != nil {
		return Courier{}, fmt.Errorf("failed to check existing courier by id: %w", err)
	}

	courier := Courier{
		ID:            dbCourier.ID,
		Name:          dbCourier.Name,
		Phone:         dbCourier.Phone,
		Status:        dbCourier.Status,
		TransportType: dbCourier.TransportType,
		CreatedAt:     dbCourier.CreatedAt,
		UpdatedAt:     dbCourier.UpdatedAt,
	}

	return courier, nil
}

func (s *Usecase) Update(ctx context.Context, c Courier) error {
	dbCourier := courier.CourierDB{
		ID:            c.ID,
		Name:          c.Name,
		Phone:         c.Phone,
		Status:        c.Status,
		TransportType: c.TransportType,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
	err := s.repo.Update(ctx, dbCourier)
	if errors.Is(err, courierErrors.ErrCourierNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to update courier: %w", err)
	}

	return nil
}

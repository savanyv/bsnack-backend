package usecase

import (
	"context"
	"time"

	dtos "github.com/savanyv/bsnack-backend/internal/dto"
	"github.com/savanyv/bsnack-backend/internal/repository"
)

type CustomerUsecase interface {
	GetCustomers(ctx context.Context, query dtos.CustomerQuery) ([]dtos.CustomerResponse, error)
}

type customerUsecase struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerUsecase(cr repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{
		customerRepo: cr,
	}
}

func (u *customerUsecase) GetCustomers(ctx context.Context, query dtos.CustomerQuery) ([]dtos.CustomerResponse, error) {
	customers, err := u.customerRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	month := query.Month
	if month.IsZero() {
		month = time.Now()
	}

	startOfMonth := time.Date(
		month.Year(),
		month.Month(),
		1, 0, 0, 0, 0,
		month.Location(),
	)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	responses := make([]dtos.CustomerResponse, 0)
	for _, c := range customers {
		isNew := c.CreatedAt.After(startOfMonth) && c.CreatedAt.Before(endOfMonth)

		responses = append(responses, dtos.CustomerResponse{
			ID: c.ID.String(),
			Name: c.Name,
			Point: c.Point,
			IsNew: isNew,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
		})
	}

	return responses, nil
}

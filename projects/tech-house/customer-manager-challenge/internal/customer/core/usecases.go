package core

import (
	"context"

	types "github.com/devpablocristo/customer-manager/pkg/types"
	domain "github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/core/domain"
	ports "github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/core/ports"
)

type UseCases struct {
	repo ports.Repository
}

func NewUseCases(r ports.Repository) ports.UseCases {
	return &UseCases{
		repo: r,
	}
}

func (uc *UseCases) GetCustomers(ctx context.Context) ([]domain.Customer, error) {
	customers, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to get customers",
			err,
		)
	}
	return customers, nil
}

func (uc *UseCases) GetCustomerByID(ctx context.Context, ID int64) (*domain.Customer, error) {
	customer, err := uc.repo.GetByID(ctx, ID)
	if err != nil {
		if types.IsNotFound(err) {
			return nil, types.NewError(
				types.ErrNotFound,
				"customer not found",
				err,
			)
		}
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to get customer",
			err,
		)
	}
	return customer, nil
}

func (uc *UseCases) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	if err := uc.repo.Create(ctx, customer); err != nil {
		if types.IsConflict(err) {
			return err // Propagamos el error de conflicto tal cual
		}
		return types.NewError(
			types.ErrOperationFailed,
			"failed to create customer",
			err,
		)
	}
	return nil
}

func (uc *UseCases) GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	customer, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if types.IsNotFound(err) {
			return nil, types.NewError(
				types.ErrNotFound,
				"customer not found",
				err,
			)
		}
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to get customer by email",
			err,
		)
	}
	return customer, nil
}

func (uc *UseCases) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
	if err := uc.repo.Update(ctx, customer); err != nil {
		if types.IsNotFound(err) {
			return err
		}
		if types.IsConflict(err) {
			return err
		}
		return types.NewError(
			types.ErrOperationFailed,
			"failed to update customer",
			err,
		)
	}
	return nil
}

func (uc *UseCases) DeleteCustomer(ctx context.Context, ID int64) error {
	if err := uc.repo.Delete(ctx, ID); err != nil {
		if types.IsNotFound(err) {
			return err
		}
		return types.NewError(
			types.ErrOperationFailed,
			"failed to delete customer",
			err,
		)
	}
	return nil
}

func (uc *UseCases) GetKPI(ctx context.Context) (*domain.KPI, error) {
	customers, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to calculate KPI",
			err,
		)
	}

	return calculateKPI(customers), nil
}

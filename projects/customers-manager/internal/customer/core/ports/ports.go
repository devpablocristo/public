package ports

import (
	"context"

	"github.com/devpablocristo/tech-house/projects/customers-manager/internal/customer/core/domain"
)

type UseCases interface {
	GetCustomers(context.Context) ([]domain.Customer, error)
	GetCustomerByID(context.Context, int64) (*domain.Customer, error)
	GetCustomerByEmail(context.Context, string) (*domain.Customer, error)
	CreateCustomer(context.Context, *domain.Customer) error
	UpdateCustomer(context.Context, *domain.Customer) error
	DeleteCustomer(context.Context, int64) error
	GetKPI(context.Context) (*domain.KPI, error)
}

type Repository interface {
	GetAll(context.Context) ([]domain.Customer, error)
	GetByID(context.Context, int64) (*domain.Customer, error)
	Create(context.Context, *domain.Customer) error
	Update(context.Context, *domain.Customer) error
	Delete(context.Context, int64) error
	GetByEmail(context.Context, string) (*domain.Customer, error)
}

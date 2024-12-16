package transport

import (
	"time"

	domain "github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/customer/core/domain"
)

// DTOs
type CustomerDataModel struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	Age       int       `db:"age"`
	BirthDate time.Time `db:"birth_date"`
}

// Mappers
func CustomerDataModelToDomain(model *CustomerDataModel) *domain.Customer {
	return &domain.Customer{
		ID:        model.ID,
		Name:      model.Name,
		LastName:  model.LastName,
		Email:     model.Email,
		Phone:     model.Phone,
		Age:       model.Age,
		BirthDate: model.BirthDate,
	}
}

func DomainToCustomerDataModel(customer *domain.Customer) *CustomerDataModel {
	return &CustomerDataModel{
		ID:        customer.ID,
		Name:      customer.Name,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Age:       customer.Age,
		BirthDate: customer.BirthDate,
	}
}

func CustomerDataModelListToDomainList(models []*CustomerDataModel) []*domain.Customer {
	customers := make([]*domain.Customer, len(models))
	for i, model := range models {
		customers[i] = CustomerDataModelToDomain(model)
	}
	return customers
}

func DomainListToCustomerDataModelList(customers []*domain.Customer) []*CustomerDataModel {
	models := make([]*CustomerDataModel, len(customers))
	for i, customer := range customers {
		models[i] = DomainToCustomerDataModel(customer)
	}
	return models
}

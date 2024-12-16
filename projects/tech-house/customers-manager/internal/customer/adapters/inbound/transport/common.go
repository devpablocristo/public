package transport

import (
	"time"

	"github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/core/domain"
)

// DTOs
type CustomerJson struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Phone     string    `json:"phone"`
	Age       int       `json:"age" binding:"required"`
	BirthDate time.Time `json:"birth_date" binding:"required"`
}

// Mappers
func CustomerJsonToDomain(c *CustomerJson) *domain.Customer {
	return &domain.Customer{
		ID:        c.ID,
		Name:      c.Name,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
		Age:       c.Age,
		BirthDate: c.BirthDate,
	}
}

func DomainToCustomerJson(customer *domain.Customer) *CustomerJson {
	return &CustomerJson{
		ID:        customer.ID,
		Name:      customer.Name,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Age:       customer.Age,
		BirthDate: customer.BirthDate,
	}
}

func CustomerJsonListToDomainList(customers []CustomerJson) []domain.Customer {
	if len(customers) == 0 {
		return []domain.Customer{}
	}

	domainCustomers := make([]domain.Customer, len(customers))
	for i, customer := range customers {
		domainCustomers[i] = *CustomerJsonToDomain(&customer)
	}
	return domainCustomers
}

func DomainListToCustomerJsonList(customers []domain.Customer) []CustomerJson {
	if len(customers) == 0 {
		return []CustomerJson{}
	}

	jsonCustomers := make([]CustomerJson, len(customers))
	for i, customer := range customers {
		jsonCustomers[i] = *DomainToCustomerJson(&customer)
	}
	return jsonCustomers
}

// Responses
type MessageResponse struct {
	Message string `json:"message"`
}

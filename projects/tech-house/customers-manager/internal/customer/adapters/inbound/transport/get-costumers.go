package transport

import (
	"github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/core/domain"
)

// Mappers

// Presenter
func GetCustomersPresenter(customers []domain.Customer) []CustomerJson {
	response := make([]CustomerJson, len(customers))
	for i, customer := range customers {
		response[i] = *DomainToCustomerJson(&customer)
	}

	return response
}

// Response
type GetCustomersResponse struct {
	Customers []CustomerJson `json:"customers"`
}

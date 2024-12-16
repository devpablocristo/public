package transport

import (
	"github.com/devpablocristo/tech-house/projects/customers-manager/internal/customer/core/domain"
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

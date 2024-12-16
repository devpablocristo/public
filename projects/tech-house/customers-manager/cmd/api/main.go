package main

import (
	"context"
	"log"

	config "github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/config"

	custin "github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/adapters/inbound"
	custout "github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/adapters/outbound"
	custcore "github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/core"
)

func init() {
	config.Load()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customerRepository, err := custout.NewRepository()
	if err != nil {
		log.Fatalf("SQLite error: %v", err)
	}

	customerUsecases := custcore.NewUseCases(customerRepository)

	customerHandler, err := custin.NewHandler(customerUsecases)
	if err != nil {
		log.Fatalf("costumer Handler error: %v", err)
	}

	err = customerHandler.Start(ctx)
	if err != nil {
		log.Fatalf("Gin Server error at start: %v", err)
	}
}

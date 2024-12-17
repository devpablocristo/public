package main

import (
	"context"
	"log"

	config "github.com/devpablocristo/tech-house/projects/customers-manager/internal/config"

	custin "github.com/devpablocristo/tech-house/projects/customers-manager/internal/customer/adapters/inbound"
	custout "github.com/devpablocristo/tech-house/projects/customers-manager/internal/customer/adapters/outbound"
	custcore "github.com/devpablocristo/tech-house/projects/customers-manager/internal/customer/core"
)

func init() {
	if err := config.Load(); err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
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

package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"

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
	customerRepository, err := custout.NewRepository()
	if err != nil {
		log.Fatalf("SQLite error: %v", err)
	}

	customerUsecases := custcore.NewUseCases(customerRepository)

	lambdaHandler, err := custin.NewLambdaHandler(customerUsecases)
	if err != nil {
		panic(err)
	}

	lambda.Start(lambdaHandler.HandleRequest)
}

package pkggorm

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/devpablocristo/customer-manager/pkg/gorm/defs"
)

var (
	instance  defs.DBClient
	once      sync.Once
	initError error
)

// gormClient es la implementación del cliente de Gorm
type gormClient struct {
	client  *gorm.DB
	address string
}

// InitializeGormClient inicializa el cliente de Gorm
func InitializeGormClient(config dbConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid Gorm configuration: %w", err)
	}

	once.Do(func() {
		client := &gormClient{}
		initError = client.connect(config)
		if initError != nil {
			instance = nil
		} else {
			instance = client
		}
	})
	return initError
}

// GetGormInstance devuelve la instancia de Gorm
func GetGormInstance() (defs.DBClient, error) {
	if instance == nil {
		return nil, fmt.Errorf("gorm client is not initialized")
	}
	return instance, nil
}

// connect conecta el cliente de Gorm
func (c *gormClient) connect(config dbConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Host, config.User, config.Password, config.DBName, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to Gorm: %w", err)
	}

	c.client = db
	c.address = config.Host

	return nil
}

// Client devuelve el cliente de Gorm
func (c *gormClient) Client() *gorm.DB {
	return c.client
}

// Address devuelve la dirección del cliente de Gorm
func (c *gormClient) Address() string {
	return c.address
}

func (c *gormClient) AutoMigrate(models ...any) error {
	return c.client.AutoMigrate(models...)
}

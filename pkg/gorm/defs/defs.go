package defs

import "gorm.io/gorm"

// DBClient define la interfaz para interactuar con el cliente de Gorm
type DBClient interface {
	Client() *gorm.DB
	Address() string
	AutoMigrate(models ...any) error // Actualizar para aceptar múltiples modelos
}

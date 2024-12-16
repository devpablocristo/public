package pkggorm

import (
	"fmt"
)

// dbConfig define la configuración para el cliente de Gorm
type dbConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

// newDBConfig crea una nueva configuración de Gorm
func newDBConfig(host, user, password, dbname string, port int) dbConfig {
	return dbConfig{
		Host:     host,
		User:     user,
		Password: password,
		DBName:   dbname,
		Port:     port,
	}
}

// Validate verifica si la configuración es válida
func (c dbConfig) Validate() error {
	if c.Host == "" || c.User == "" || c.Password == "" || c.DBName == "" || c.Port == 0 {
		return fmt.Errorf("incomplete Gorm configuration")
	}
	return nil
}

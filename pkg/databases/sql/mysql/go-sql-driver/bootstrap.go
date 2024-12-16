package pkgmysql

import (
	"github.com/spf13/viper"

	"github.com/devpablocristo/customer-manager/pkg/databases/sql/mysql/go-sql-driver/defs"
)

// Bootstrap inicializa la configuración y crea una instancia de repositorio MySQL.
func Bootstrap() (defs.Repository, error) {
	config := config{
		User:     viper.GetString("MYSQL_USER"),
		Password: viper.GetString("MYSQL_PASSWORD"),
		Host:     viper.GetString("MYSQL_HOST"),
		Port:     viper.GetString("MYSQL_PORT"),
		Database: viper.GetString("MYSQL_DATABASE"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

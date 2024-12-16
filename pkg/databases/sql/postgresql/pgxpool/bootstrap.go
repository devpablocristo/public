package pkgpostgresql

import (
	"github.com/spf13/viper"

	"github.com/devpablocristo/customer-manager/pkg/databases/sql/postgresql/pgxpool/defs"
)

// NOTE: Diseñado para establer conexion con 1 base de datos durante la ejecución de la app
func Bootstrap() (defs.Repository, error) {
	config := newConfig(
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_HOST"),
		viper.GetString("POSTGRES_PORT"),
		viper.GetString("POSTGRES_MIGRATIONS_DIR"),
		viper.GetString("POSTGRES_DB"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

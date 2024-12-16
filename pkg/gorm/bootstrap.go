package pkggorm

import (
	"github.com/spf13/viper"

	"github.com/devpablocristo/customer-manager/pkg/gorm/defs"
)

func Bootstrap() (defs.DBClient, error) {
	config := newDBConfig(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetInt("DB_PORT"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	if err := InitializeGormClient(config); err != nil {
		return nil, err
	}

	return GetGormInstance()
}

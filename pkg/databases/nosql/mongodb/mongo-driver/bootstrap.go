package mongodbdriver

import (
	"github.com/spf13/viper"

	defs "github.com/devpablocristo/customer-manager/pkg/databases/nosql/mongodb/mongo-driver/defs"
)

func Bootstrap() (defs.Repository, error) {
	config := newConfig(
		viper.GetString("MONGO_USER"),
		viper.GetString("MONGO_PASSWORD"),
		viper.GetString("MONGO_HOST"),
		viper.GetString("MONGO_PORT"),
		viper.GetString("MONGO_DATABASE"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

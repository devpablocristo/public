package pkgredis

import (
	"github.com/spf13/viper"

	defs "github.com/devpablocristo/customer-manager/pkg/cache/redis/v8/defs"
)

func Bootstrap() (defs.Cache, error) {
	config := newConfig(
		viper.GetString("REDIS_ADDRESS"),
		viper.GetString("REDIS_PASSWORD"),
		viper.GetInt("REDIS_DB"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newCache(config)
}

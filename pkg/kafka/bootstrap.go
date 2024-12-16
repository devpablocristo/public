package pkgafka

import (
	"github.com/spf13/viper"

	"github.com/devpablocristo/customer-manager/pkg/kafka/defs"
)

func Bootstrap(brokersKey, groupIDKey string) (defs.Service, error) {
	config := newConfig(
		viper.GetStringSlice(brokersKey),
		viper.GetString(groupIDKey),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newService(config)
}

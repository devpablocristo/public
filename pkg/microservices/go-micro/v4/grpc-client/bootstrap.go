package pkggomicro

import (
	"github.com/spf13/viper"

	defs "github.com/devpablocristo/customer-manager/pkg/microservices/go-micro/v4/grpc-client/defs"
)

func Bootstrap() (defs.Client, error) {
	config := newConfig(
		viper.GetString("CONSUL_ADDRESS"),
		viper.GetString("GRPC_SERVER_NAME"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newClient(config)
}

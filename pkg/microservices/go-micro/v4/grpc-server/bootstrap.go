package pkggomicro

import (
	"github.com/spf13/viper"

	defs "github.com/devpablocristo/customer-manager/pkg/microservices/go-micro/v4/grpc-server/defs"
)

func Bootstrap() (defs.Server, error) {
	config := newConfig(
		viper.GetString("GRPC_SERVER_NAME"),
		viper.GetString("GRPC_SERVER_HOST"),
		viper.GetInt("GRPC_SERVER_PORT"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newServer(config)
}

package pkggomicro

import (
	"github.com/spf13/viper"

	pkgclient "github.com/devpablocristo/customer-manager/pkg/microservices/go-micro/v4/grpc-client/defs"
	pkgserver "github.com/devpablocristo/customer-manager/pkg/microservices/go-micro/v4/grpc-server/defs"
	pkgservice "github.com/devpablocristo/customer-manager/pkg/microservices/go-micro/v4/micro-service/defs"
	pkgbroker "github.com/devpablocristo/customer-manager/pkg/microservices/go-micro/v4/rabbitmq-broker/defs"
)

func Bootstrap(server pkgserver.Server, client pkgclient.Client, broker pkgbroker.Broker) (pkgservice.Service, error) {
	config := newConfig(
		server.GetServer(),
		client.GetClient(),
		broker.GetBroker(),
		viper.GetString("CONSUL_ADDRESS"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newService(config)
}

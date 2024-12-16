package pkggrpcserver

import (
	"github.com/spf13/viper"

	defs "github.com/devpablocristo/customer-manager/pkg/grpc/server/defs"
)

// Bootstrap inicializa y devuelve una instancia de servidor gRPC
func Bootstrap() (defs.Server, error) {
	host := viper.GetString("GRPC_SERVER_HOST")
	if host == "" {
		host = "0.0.0.0" // Valor predeterminado si no se especifica
	}

	config := newConfig(
		host, // viper.GetString("GRPC_SERVER_HOST"), // si es necesario
		viper.GetInt("GRPC_SERVER_PORT"),
		nil, // Configuración TLS, si es necesario
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newServer(config)
}

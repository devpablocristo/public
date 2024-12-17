package pkgaws

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/devpablocristo/tech-house/pkg/aws/defs"
)

// Bootstrap inicializa y retorna un Stack AWS basado en la configuración del entorno
func Bootstrap() (defs.Stack, error) {
	// Validar y obtener el provider
	provider := viper.GetString("AWS_PROVIDER")
	if provider == "" {
		return nil, fmt.Errorf("AWS_PROVIDER is required (aws or localstack)")
	}
	if provider != defs.ProviderAWS && provider != defs.ProviderLocalstack {
		return nil, fmt.Errorf("invalid AWS_PROVIDER: %s", provider)
	}

	// Validar credenciales requeridas
	accessKey := viper.GetString("AWS_ACCESS_KEY_ID")
	secretKey := viper.GetString("AWS_SECRET_ACCESS_KEY")
	region := viper.GetString("AWS_REGION")

	// Crear configuración base
	opts := []ConfigOption{
		WithCredentials(accessKey, secretKey),
		WithRegion(region),
	}

	// Validar y configurar servicios si están especificados
	if servicesStr := viper.GetString("AWS_SERVICES"); servicesStr != "" {
		services := strings.Split(servicesStr, ",")
		for _, service := range services {
			if !defs.ValidServices[strings.TrimSpace(service)] {
				return nil, fmt.Errorf("invalid service: %s", service)
			}
		}
		opts = append(opts, WithServices(services))
	}

	// Configuración específica de Localstack
	if provider == defs.ProviderLocalstack {
		endpoint := viper.GetString("AWS_LOCALSTACK_ENDPOINT")
		if endpoint == "" {
			return nil, fmt.Errorf("AWS_LOCALSTACK_ENDPOINT is required for localstack")
		}

		opts = append(opts, WithLocalstackConfig(
			endpoint,
			viper.GetInt("AWS_EDGE_PORT"),
			viper.GetInt("AWS_WEB_UI_PORT"),
		))
	}

	// Crear y validar la configuración
	config := NewConfig(provider, opts...)
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Crear el stack usando el factory
	factory, err := NewStackFactory(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create stack factory: %w", err)
	}

	return factory.CreateStack(config)
}

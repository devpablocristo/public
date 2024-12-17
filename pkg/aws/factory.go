package pkgaws

import (
	"fmt"

	defs "github.com/devpablocristo/tech-house/pkg/aws/defs"
	localstack "github.com/devpablocristo/tech-house/pkg/aws/localstack"
	realaws "github.com/devpablocristo/tech-house/pkg/aws/realstack"
)

// StackFactory define la interfaz para la creación de stacks
type StackFactory interface {
	// CreateStack crea un nuevo stack basado en la configuración proporcionada
	CreateStack(config defs.Config) (defs.Stack, error)
}

// awsProvider implementa StackFactory para AWS real
type awsProvider struct{}

// localstackProvider implementa StackFactory para Localstack
type localstackProvider struct{}

// NewStackFactory crea un nuevo factory basado en el provider especificado
func NewStackFactory(provider string) (StackFactory, error) {
	// Validar que el provider sea válido
	if provider == "" {
		return nil, &ConfigError{
			Field:   "provider",
			Message: "provider cannot be empty",
		}
	}

	switch provider {
	case defs.ProviderAWS:
		return &awsProvider{}, nil
	case defs.ProviderLocalstack:
		return &localstackProvider{}, nil
	default:
		return nil, &ConfigError{
			Field:   "provider",
			Message: fmt.Sprintf("unsupported provider: %s", provider),
		}
	}
}

// CreateStack implementación para AWS real
func (p *awsProvider) CreateStack(config defs.Config) (defs.Stack, error) {
	// Validar que la configuración coincida con el provider
	if config.GetProvider() != defs.ProviderAWS {
		return nil, &ConfigError{
			Field:   "provider",
			Message: "config provider does not match AWS provider",
		}
	}

	stack, err := realaws.NewStack(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS stack: %w", err)
	}

	return stack, nil
}

// CreateStack implementación para Localstack
func (p *localstackProvider) CreateStack(config defs.Config) (defs.Stack, error) {
	// Validar que la configuración coincida con el provider
	if config.GetProvider() != defs.ProviderLocalstack {
		return nil, &ConfigError{
			Field:   "provider",
			Message: "config provider does not match Localstack provider",
		}
	}

	// Validar endpoint para Localstack
	if config.GetEndpoint() == "" {
		return nil, &ConfigError{
			Field:   "endpoint",
			Message: "endpoint is required for Localstack",
		}
	}

	stack, err := localstack.NewStack(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Localstack stack: %w", err)
	}

	return stack, nil
}

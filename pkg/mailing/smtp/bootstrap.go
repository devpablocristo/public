package pkgsmtp

import (
	"fmt"

	"github.com/spf13/viper"

	defs "github.com/devpablocristo/customer-manager/pkg/mailing/smtp/defs"
)

func Bootstrap() (defs.Service, error) {
	config := newConfig(
		viper.GetString("SMTP_HOST"),
		viper.GetString("SMTP_PORT"),
		viper.GetString("SMTP_FROM"),
		viper.GetString("SMTP_USERNAME"),
		viper.GetString("SMTP_PASSWORD"),
		viper.GetString("SMTP_IDENTITY"),
		viper.GetString("VERIFICATION_URL"),
	)

	// Validar la configuración
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("SMTP config error: %w", err)
	}

	// Crear el servicio SMTP con la configuración
	return newService(config)
}

package pkghcl

import (
	"github.com/spf13/viper"

	"github.com/devpablocristo/golang-monorepo/pkg/rest/clients/net-http/defs"
)

func Bootstrap(tokenEndPoint, clientID, clientSecret string, additionalParams map[string]string) (defs.Client, defs.Config, error) {
	if tokenEndPoint == "" {
		tokenEndPoint = viper.GetString("HTTP_CLIENT_ENDPOINT_KEY")
	}
	if clientID == "" {
		tokenEndPoint = viper.GetString("HTTP_CLIENT_CLIENT_ID")
	}
	if clientSecret == "" {
		clientSecret = viper.GetString("HTTP_CLIENT_SECRET")
	}
	if additionalParams == nil {
		additionalParams = viper.GetStringMapString("HTTP_CLIENT_ADD_PARAMS")
	}

	config := newConfig(
		tokenEndPoint,
		clientID,
		clientSecret,
		additionalParams,
	)

	if err := config.Validate(); err != nil {
		return nil, nil, err
	}

	client, err := newClient(config)
	if err != nil {
		return nil, nil, err
	}

	// NOTE: realmente necesito enviar config???
	return client, config, nil
}

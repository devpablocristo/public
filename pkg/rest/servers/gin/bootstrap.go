package pkggin

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	defs "github.com/devpablocristo/tech-house/pkg/rest/servers/gin/defs"
)

func Bootstrap(isTest bool) (defs.Server, error) {
	if gin.Mode() == gin.TestMode {
		return newTestServer()
	}

	config := newConfig(
		viper.GetString("WEB_SERVER_PORT"),
		viper.GetString("API_VERSION"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newServer(config)
}

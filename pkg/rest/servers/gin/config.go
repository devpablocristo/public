package pkggin

import (
	"fmt"

	defs "github.com/devpablocristo/golang-monorepo/pkg/rest/servers/gin/defs"
)

type config struct {
	routerPort string
	apiVersion string
}

func newConfig(routerPort, ApiVersion string) defs.Config {
	return &config{
		routerPort: routerPort,
		apiVersion: ApiVersion,
	}
}

func (c *config) GetRouterPort() string {
	return c.routerPort
}

func (c *config) SetRouterPort(routerPort string) {
	c.routerPort = routerPort
}

func (c *config) GetApiVersion() string {
	return c.apiVersion
}

func (c *config) SetApiVersion(ApiVersion string) {
	c.apiVersion = ApiVersion
}

func (c *config) Validate() error {
	if c.routerPort == "" {
		return fmt.Errorf("router port is not configured")
	}
	return nil
}

package pkggin

import (
	defs "github.com/devpablocristo/golang-monorepo/pkg/rest/servers/gin/defs"
	"github.com/gin-gonic/gin"
)

// NewTestConfig crea una configuración específica para tests
func NewTestConfig() defs.Config {
	cfg := newConfig("", "")
	cfg.SetRouterPort("8080")
	cfg.SetApiVersion("v1")
	return cfg
}

// NewTestServer crea un servidor configurado para tests
func NewTestServer() (defs.Server, error) {
	gin.SetMode(gin.TestMode)
	cfg := NewTestConfig()

	// No usamos singleton para tests
	r := gin.New()
	return &server{
		config: cfg,
		router: r,
	}, nil
}

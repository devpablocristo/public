package pkggin

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	defs "github.com/devpablocristo/tech-house/pkg/rest/servers/gin/defs"
)

var (
	instance  defs.Server
	once      sync.Once
	initError error
)

type server struct {
	router *gin.Engine
	config defs.Config
}

func newServer(config defs.Config) (defs.Server, error) {
	once.Do(func() {
		err := config.Validate()
		if err != nil {
			initError = err
			return
		}

		r := gin.Default()
		instance = &server{
			config: config,
			router: r,
		}
	})
	return instance, initError
}

func newTestServer() (defs.Server, error) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Crear una configuración de prueba
	testConfig := &config{
		routerPort: "8080",
		apiVersion: "v1",
	}

	// Crear una nueva instancia del servidor para tests
	// No usamos el singleton pattern en tests
	return &server{
		router: r,
		config: testConfig,
	}, nil
}

func (server *server) RunServer(ctx context.Context) error {
	return server.router.Run(":" + server.config.GetRouterPort())
}

func (server *server) GetRouter() *gin.Engine {
	return server.router
}

func (server *server) SetRouter() *gin.Engine {
	return server.router
}

func (server *server) GetApiVersion() string {
	return server.config.GetApiVersion()
}

// WrapH envuelve un http.Handler en un gin.HandlerFunc.
func (server *server) WrapH(h http.Handler) gin.HandlerFunc {
	return gin.WrapH(h)
}

package pkgswagger

import (
	"encoding/json"
	"net/http"
	"path"
	"sync"

	"github.com/devpablocristo/tech-house/pkg/swagger/defs"
)

var (
	instance  defs.Service
	once      sync.Once
	initError error
)

type service struct {
	config defs.Config
}

func newService(config defs.Config) (defs.Service, error) {
	once.Do(func() {
		err := config.Validate()
		if err != nil {
			initError = err
			return
		}
		instance = &service{
			config: config,
		}
	})
	return instance, initError
}

func (s *service) Setup(addRoute func(defs.HandlerConfig)) error {
	if !s.config.IsEnabled() {
		return nil
	}

	// Configurar rutas base de Swagger
	addRoute(defs.HandlerConfig{
		Path:   "/swagger.json",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			s.serveSwaggerSpec(w, r)
		},
	})

	addRoute(defs.HandlerConfig{
		Path:   "/swagger/",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			s.serveSwaggerUI(w, r)
		},
	})

	return nil
}

func (s *service) serveSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	spec := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       s.config.GetTitle(),
			"description": s.config.GetDescription(),
			"version":     s.config.GetVersion(),
		},
		"host":     s.config.GetHost(),
		"basePath": s.config.GetBasePath(),
		"schemes":  s.config.GetSchemes(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spec)
}

func (s *service) serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Si es la ruta base de swagger, redirigir al index
	if r.URL.Path == "/swagger" || r.URL.Path == "/swagger/" {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
		return
	}

	// Servir archivos estáticos de Swagger UI
	filePath := path.Clean(r.URL.Path[len("/swagger/"):])
	http.ServeFile(w, r, "swagger-ui/"+filePath)
}

func (s *service) GetConfig() defs.Config {
	return s.config
}

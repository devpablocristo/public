package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	initconf "github.com/devpablocristo/tech-house/pkg/config/init-config"
	mwr "github.com/devpablocristo/tech-house/pkg/rest/middlewares/gin"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	auth mwr.Config
}

func Load() error {
	var loadErr error
	once.Do(func() {
		initconf.LoadConfig("config/.env")

		secretKey, err := getEnv("JWT_SECRET_KEY")
		if err != nil {
			loadErr = err
			return
		}

		cfg = &Config{
			auth: mwr.Config{
				SecretKey:   secretKey,
				TokenLookup: "header:Authorization",
				TokenPrefix: "Bearer ",
			},
		}
	})
	return loadErr
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}

// Auth returns middleware auth configuration
func Auth() mwr.Config {
	if cfg == nil {
		log.Fatal("configuration not loaded")
	}
	return cfg.auth
}

// MustLoad loads the configuration or panics
func MustLoad() {
	if err := Load(); err != nil {
		log.Fatal(err)
	}
}

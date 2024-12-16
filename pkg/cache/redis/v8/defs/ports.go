package defs

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Client() *redis.Client
	Close()
	Set(context.Context, string, any, time.Duration) error
	Get(context.Context, string) (string, error)
}

// Config define el puerto para la configuración de Redis
type Config interface {
	GetAddress() string
	GetPassword() string
	GetDB() int
	Validate() error
}

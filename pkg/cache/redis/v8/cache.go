package pkgredis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	defs "github.com/devpablocristo/customer-manager/pkg/cache/redis/v8/defs"
)

var (
	instance  defs.Cache
	once      sync.Once
	initError error
)

type Cache struct {
	client *redis.Client
}

// newCache inicializa la instancia del cache Redis utilizando el patrón singleton
func newCache(c defs.Config) (defs.Cache, error) {
	once.Do(func() {
		client := &Cache{}
		initError = client.connect(c)
		if initError != nil {
			instance = nil
		} else {
			instance = client
		}
	})
	return instance, initError
}

// connect conecta al servidor Redis utilizando los getters de la configuración
func (ch *Cache) connect(c defs.Config) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.GetAddress(),
		Password: c.GetPassword(),
		DB:       c.GetDB(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	ch.client = rdb
	return nil
}

// Set almacena un valor en Redis con una clave y un tiempo de expiración
func (ch *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return ch.client.Set(ctx, key, value, expiration).Err()
}

// Get recupera un valor de Redis usando una clave
func (ch *Cache) Get(ctx context.Context, key string) (string, error) {
	result, err := ch.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}
	return result, nil
}

// Close cierra la conexión con el servidor Redis
func (ch *Cache) Close() {
	if ch.client != nil {
		ch.client.Close()
	}
}

// Client devuelve el cliente Redis
func (ch *Cache) Client() *redis.Client {
	return ch.client
}

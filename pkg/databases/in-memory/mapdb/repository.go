package pkgmapdb

import (
	"sync"

	defs "github.com/devpablocristo/customer-manager/pkg/databases/in-memory/mapdb/defs"
)

var (
	instance defs.Repository
	once     sync.Once
)

type service struct {
	db map[string]any
}

func newRepository() defs.Repository {
	once.Do(func() {
		instance = &service{
			db: make(map[string]any),
		}
	})
	return instance
}

func (c *service) GetDb() map[string]any {
	return c.db
}

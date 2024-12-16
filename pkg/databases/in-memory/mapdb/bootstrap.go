package pkgmapdb

import (
	defs "github.com/devpablocristo/customer-manager/pkg/databases/in-memory/mapdb/defs"
)

func Boostrap() defs.Repository {
	return newRepository()
}

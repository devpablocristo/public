package pkgsqlite

import (
	"github.com/spf13/viper"

	"github.com/devpablocristo/customer-manager/pkg/databases/sql/sqlite/defs"
)

func Bootstrap() (defs.Repository, error) {
	config := newConfig(
		viper.GetString("SQLITE_DB_PATH"),
		viper.GetBool("SQLITE_IN_MEMORY"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

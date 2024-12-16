package pkgcassandra

import (
	"github.com/spf13/viper"

	defs "github.com/devpablocristo/customer-manager/pkg/databases/nosql/cassandra/gocql/defs"
)

func Bootstrap() (defs.Repository, error) {
	config := newConfig(
		viper.GetStringSlice("CASSANDRA_HOSTS"),
		viper.GetString("CASSANDRA_KEYSPACE"),
		viper.GetString("CASSANDRA_USERNAME"),
		viper.GetString("CASSANDRA_PASSWORD"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

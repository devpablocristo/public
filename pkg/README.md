# NOTE: cambiar todos los bootstraps: reciben la variable y no la el key de la env var, con la adhesion de config.go al template del servicio, esta configuracion es mejor.
ejemplo de nuevo implementaicon

Las configs se hace por .env, como es no meno repo ej .env.user
```go
func Bootstrap(dbName string) (defs.Repository, error) {
	// Si dbName no se proporciona, lee de .env
	if dbName == "" {
		dbName = viper.GetString("POSTGRES_DB_NAME")
	}

	config := newConfig(
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_HOST"),
		viper.GetString("POSTGRES_PORT"),
		viper.GetString("POSTGRES_MIGRATIONS_DIR"),
		dbName,
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}
```



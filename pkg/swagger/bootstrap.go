package pkgswagger

import (
    "os"
    "strings"

    "github.com/devpablocristo/tech-house/pkg/swagger/defs"
)

func Bootstrap() (defs.Service, error) {
    config := newConfig(
        os.Getenv("SWAGGER_TITLE"),
        os.Getenv("SWAGGER_DESCRIPTION"),
        os.Getenv("SWAGGER_VERSION"),
        os.Getenv("SWAGGER_HOST"),
        os.Getenv("SWAGGER_BASE_PATH"),
        strings.Split(os.Getenv("SWAGGER_SCHEMES"), ","),
        os.Getenv("SWAGGER_ENABLED") == "true",
    )

    if err := config.Validate(); err != nil {
        return nil, err
    }

    return newService(config)
}
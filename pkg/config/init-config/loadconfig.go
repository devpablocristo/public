package pkgconf

import (
	"log"

	pkgenvs "github.com/devpablocristo/golang-monorepo/pkg/config/godotenv"
	pkgviper "github.com/devpablocristo/golang-monorepo/pkg/config/viper"
	pkgutils "github.com/devpablocristo/golang-monorepo/pkg/utils"
)

func LoadConfig(filePaths ...string) {
	if len(filePaths) == 0 {
		log.Fatal("Fatal error: no configuration file paths provided")
	}

	foundFiles, err := pkgutils.FilesFinder(filePaths...)
	if err != nil {
		log.Fatalf("Fatal error: failed to find configuration files: %v", err)
	}

	if err := pkgenvs.LoadConfig(foundFiles...); err != nil {
		log.Fatalf("Fatal error: failed to initialize environment configuration: %v", err)
	}

	if err := pkgviper.LoadConfig(foundFiles...); err != nil {
		log.Fatalf("Fatal error: failed to initialize Viper configuration: %v", err)
	}
}

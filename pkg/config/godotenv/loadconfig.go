package pkgenvs

import (
	"errors"

	"github.com/joho/godotenv"
)

// NOTE: Rediseñar los prints
// LoadConfig loads multiple .env files using godotenv
func LoadConfig(filePaths ...string) error {
	if len(filePaths) == 0 {
		return errors.New("no environment file paths provided")
	}

	successfullyLoaded := false

	for _, filePath := range filePaths {
		if err := godotenv.Load(filePath); err != nil {
			//fmt.Printf("Error loading environment file: '%s'\n", filePath)
			continue
		}
		//fmt.Printf("Environment file loaded successfully: %s\n", filePath)
		successfullyLoaded = true
	}

	if !successfullyLoaded {
		return errors.New("failed to load any environment files")
	}

	return nil
}

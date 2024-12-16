package pkgviper

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	pkgutils "github.com/devpablocristo/customer-manager/pkg/utils"
)

// NOTE: revisar los prints y errores
func LoadConfig(filePaths ...string) error {
	if len(filePaths) == 0 {
		return errors.New("no file paths provided")
	}

	configureViper()

	loadedDirs := make(map[string]bool)
	//successfullyLoaded := false // Variable to track successful loads

	for _, configFilePath := range filePaths {
		// Try loading each file with Viper
		if err := loadViperConfig(configFilePath, loadedDirs); err != nil {
			fmt.Printf("%v\n", err) // Print error but continue
		} else {
			//successfullyLoaded = true // At least one file loaded successfully
		}
	}

	// If no file was successfully loaded, return an error
	// if !successfullyLoaded {
	// 	fmt.Println("pkgviper: WARNING: no configuration files were successfully loaded")
	// }

	return nil
}

// configureViper sets up Viper to load environment variables
func configureViper() {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//fmt.Println("pkgviper: Set to load environment variables (if any are present)")
}

// loadViperConfig loads a configuration file compatible with Viper
func loadViperConfig(configFilePath string, loadedDirs map[string]bool) error {
	fileNameWithoutExt, fileExtension, err := pkgutils.FileNameAndExtension(configFilePath)
	if err != nil {
		return fmt.Errorf("Skipping file '%s': %v", configFilePath, err)
	}

	viper.SetConfigName(fileNameWithoutExt)
	viper.SetConfigType(fileExtension)

	dir := filepath.Dir(configFilePath)
	if !loadedDirs[dir] {
		viper.AddConfigPath(dir)
		loadedDirs[dir] = true
	}

	// if err := viper.ReadInConfig(); err != nil {
	// 	return fmt.Errorf("Failed to load configuration file: '%s'", configFilePath)
	// }

	//fmt.Printf("pkgviper: Configuration file successfully loaded from: %s\n", viper.ConfigFileUsed())
	return nil
}

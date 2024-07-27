package yamlconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

func ReadConfig(subDirName string) (*any, error) {
	configDir, err := getConfigDir()
	if err != nil {
		fmt.Printf("Error determining config directory: %v\n", err)
		return nil, fmt.Errorf("TODO")
	}

	configFilePath := filepath.Join(configDir, subDirName, "config.yaml")
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return nil, err
	}
	defer file.Close()

	return nil, nil // TODO:
}

func WriteConfig(subDirName string, config *any) {
	configDir, err := getConfigDir()
	if err != nil {
		fmt.Printf("Error determining config directory: %v\n", err)
		return
	}

	configFilePath := filepath.Join(configDir, subDirName, "config.yaml")

	// Create or open the YAML file
	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(*config)
	if err != nil {
		fmt.Printf("Error encoding config content: %v\n", err)
	} else {
		fmt.Printf("Config saved successfully to %s.\n", configFilePath)
	}
}

func getConfigDir() (string, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".config")
		}
		return filepath.Abs(configDir)
	case "windows":
		return filepath.Join(os.Getenv("APPDATA")), nil
	default:
		return "", fmt.Errorf("unsupported operating system")
	}
}

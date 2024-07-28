package yamlconfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

func Read(config any, subDirName string) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configSubDirPath := filepath.Join(configDir, subDirName)
	os.MkdirAll(configSubDirPath, os.ModePerm)

	configFilePath := filepath.Join(configSubDirPath, "config.yaml")

	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlData, config); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return nil
}

func Write(subDirName string, config any) {
	configDir, err := getConfigDir()
	if err != nil {
		fmt.Printf("Error determining config directory: %v\n", err)
		return
	}

	configSubDirPath := filepath.Join(configDir, subDirName)
	os.MkdirAll(configSubDirPath, os.ModePerm)

	configFilePath := filepath.Join(configSubDirPath, "config.yaml")

	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(config)
	if err != nil {
		fmt.Printf("Error encoding config content: %v\n", err)
	} else {
		fmt.Printf("Config saved successfully to %s.\n", configFilePath)
	}
}

func GetConfigFilePath(subDirName string) string {
	configDir, _ := getConfigDir()
	configFilePath := filepath.Join(configDir, subDirName, "config.yaml")
	return configFilePath
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

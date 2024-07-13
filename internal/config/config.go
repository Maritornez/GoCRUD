package config

import (
	"os"

	"gopkg.in/yaml.v2" // Parse yaml file
)

const YamlPath = "../configs/config.yml"

// Config struct for webapp config
type Config struct {
	Database struct {
		IpAddress string `yaml:"ip_address"`
		Port      string `yaml:"port"`
		User      string `yaml:"user"`
		Pass      string `yaml:"pass"`
		DBName    string `yaml:"db_name"`
	} `yaml:"database"`
	Server struct {
		IpAddress string `yaml:"ip_address"`
		Port      string `yaml:"port"`
	} `yaml:"server"`
}

func NewConfig(path string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

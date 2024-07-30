package config

import (
	"fmt"
	"os"
	//"gopkg.in/yaml.v2" // Parse yaml file
)

// const YamlPath = "../configs/config.yml"

// type Config struct {
// 	Database struct {
// 		IpAddress string `yaml:"ip_address"`
// 		Port      string `yaml:"port"`
// 		User      string `yaml:"user"`
// 		Pass      string `yaml:"pass"`
// 		DBName    string `yaml:"db_name"`
// 	} `yaml:"database"`
// 	Server struct {
// 		IpAddress string `yaml:"ip_address"`
// 		Port      string `yaml:"port"`
// 	} `yaml:"server"`
// }

// func NewConfig(path string) (*Config, error) {
// 	// Create config structure
// 	config := &Config{}

// 	// Open config file
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	// Init new YAML decode
// 	d := yaml.NewDecoder(file)

// 	// Start YAML decoding from file
// 	if err := d.Decode(&config); err != nil {
// 		return nil, err
// 	}

// 	return config, nil
// }

type Config struct {
	Server struct {
		IpAddress string
		Port      string
	}
	Database struct {
		IpAddress string
		Port      string
		User      string
		Pass      string
		DBName    string
	}
}

func NewConfig() (*Config, error) {
	config := &Config{}

	config.Database.IpAddress = os.Getenv("DB_IP_ADDRESS")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Pass = os.Getenv("DB_PASS")
	config.Database.DBName = os.Getenv("DB_NAME")

	config.Server.IpAddress = os.Getenv("SERVER_IP_ADDRESS")
	config.Server.Port = os.Getenv("SERVER_PORT")

	// Проверим, установлены ли все переменные окружения
	if config.Database.IpAddress == "" || config.Database.Port == "" ||
		config.Database.User == "" || config.Database.Pass == "" ||
		config.Database.DBName == "" || config.Server.IpAddress == "" ||
		config.Server.Port == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return config, nil
}

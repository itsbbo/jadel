package app

import (
	"io"
	"os"

	"sigs.k8s.io/yaml"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"db"`
}

func InitConfig(filepath string) (Config, error) {
	var (
		config Config
	)

	_, err := os.Stat(filepath)
	if err != nil {
		return config, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return config, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

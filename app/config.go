package app

import (
	"io"
	"os"

	"github.com/samber/oops"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port  int  `yaml:"port"`
		Debug bool `yaml:"debug"`
	} `yaml:"server"`

	DB struct {
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		Database       string `yaml:"database"`
		SSLMode        string `yaml:"sslmode"`
		DSN            string `yaml:"datasource"`
		ConnectTimeout int    `yaml:"connect_timeout"`
	} `yaml:"db"`

	Cookie struct {
		Secure   bool   `yaml:"secure"`
		HTTPOnly bool   `yaml:"http_only"`
		SameSite string `yaml:"same_site"`
		Path     string `yaml:"path"`
		Domain   string `yaml:"domain"`
	} `yaml:"cookie"`
}

func InitConfig(filepath string) (Config, error) {
	var (
		config Config
	)

	_, err := os.Stat(filepath)
	if err != nil {
		return config, oops.In("InitConfig").Errorf("cannot stat config file: %w", err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return config, oops.In("InitConfig").Errorf("cannot open config file: %w", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return config, oops.In("InitConfig").Errorf("cannot read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, oops.In("InitConfig").Errorf("cannot unmarshal config file: %w", err)
	}

	return config, nil
}

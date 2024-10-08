package config

import (
	"os"
	"gopkg.in/yaml.v3"
)
type Config struct {
	KafkaConfig      KafkaConfig      `yaml:"kafka"`
	PostgreSQLConfig PostgreSQLConfig `yaml:"postgresql"`
}

type KafkaConfig struct {
	BootstrapServer string `yaml:"bootstrap-server"`
	Topic           string `yaml:"topic"`
}
type PostgreSQLConfig struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	MaxConnections int32  `yaml:"max-connections"`
	MinConnections int32  `yaml:"min-connections"`
	DatabaseType   string `yaml:"database-type"`
	Database       string `yaml:"db"`
	DbInit         string `yaml:"dbinit"`
}

func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
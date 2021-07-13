package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Config struct {
	Mongo        MongoConfig
	HTTP         HTTPConfig         `yaml:"http"`
	TokenManager TokenManagerConfig `yaml:"tokenManager"`

	Salt string
}

type MongoConfig struct {
	URI      string
	Username string
	Password string
	Name     string
}

type HTTPConfig struct {
	Port               string        `yaml:"port"`
	ReadTimeout        time.Duration `yaml:"readTimeout"`
	WriteTimeout       time.Duration `yaml:"writeTimeout"`
	MaxHeaderMegabytes int           `yaml:"maxHeaderMB"`
}

type TokenManagerConfig struct {
	SigningKey string
	TokenTTL   time.Duration `yaml:"tokenTTL"`
}

func NewConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := envconfig.Process("mongo", &cfg.Mongo); err != nil {
		return nil, err
	}

	if err := envconfig.Process("token", &cfg.TokenManager); err != nil {
		return nil, err
	}

	if err := envconfig.Process("hasher", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

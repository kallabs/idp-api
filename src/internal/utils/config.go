package utils

import (
	"errors"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	ServerAddress string `yaml:"server_address" envconfig:"SERVER_ADDRESS"`
	SecretKey     string `yaml:"secret_key" envconfig:"SECRET_KEY"`
	RsaPublicKey  string `yaml:"rsa_public_key" envconfig:"RSA_PUBLIC_KEY"`
	RsaPrivateKey string `yaml:"rsa_private_key" envconfig:"RSA_PRIVATE_KEY"`
	DatabaseUri   string `yaml:"database_uri" envconfig:"DATABASE_URI"`
	Jwt           struct {
		Algorithm       string   `yaml:"algorithm"`
		Issuer          string   `yaml:"issuer"`
		Audience        []string `yaml:"audience"`
		AccessLifetime  int      `yaml:"access_lifetime"`
		RefreshLifetime int      `yaml:"refresh_lifetime"`
	}
}

var Conf = &AppConfig{}

// Initializes AppConfig with environment variables
func (config *AppConfig) ReadEnv() {
	envconfig.Process("app", config)
}

// Initializes AppConfig with data in yaml config file
func (config *AppConfig) ReadYaml(filepath string) error {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return errors.New("error reading config file")
	}

	yaml.Unmarshal(yamlFile, config)

	return nil
}

func LoadConfig() error {
	// Load default config file
	if err := Conf.ReadYaml("/usr/src/config.yaml"); err != nil {
		return err
	}
	// Load specific configuration with env variables
	Conf.ReadEnv()

	return nil
}

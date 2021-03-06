package config

import (
	"github.com/saegewerk/pufferApi/pkg/cache"
	"github.com/saegewerk/pufferApi/pkg/service"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Services map[string]service.Config
	Port     string
	Cache    cache.Config
}

func YAML() (config *Config, err error) {
	return YAMLfromFile("./Pufferfile.yml")
}

func YAMLfromFile(file string) (config *Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return &Config{}, err
	}
	err = yaml.NewDecoder(f).Decode(&config)
	return config, err
}

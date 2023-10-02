package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

func ParseAndValidate(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

package config

import (
	"os"

	"github.com/BurntSushi/toml"

	"github.com/keepcalmist/chat-service/internal/validator"
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

	err = validator.Validator.Struct(cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

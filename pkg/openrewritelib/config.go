package openrewritelib

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Recipes      []string `yaml:"recipes"`
	Styles       []string `yaml:"styles"`
	ExcludeFiles []string `yaml:"excludeFiles"`
}

func LoadProjectConfig(file string) (Config, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read %s: %w", file, err)
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse %s: %w", file, err)
	}

	return config, nil
}

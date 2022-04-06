package cmd

import (
	_ "embed"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

//go:embed config.yml
var yamlConfig []byte
var cfg Config

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Version string `yaml:"version"`
	Repo    string `yaml:"repo"`
}

func loadConfig() (Config, error) {
	err := yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("loadConfig: %w", err)
	}

	return cfg, nil
}

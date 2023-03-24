package config

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"sqser/models"
)

type Config struct {
	Config *Model
}

//go:embed config.yaml
var config []byte

type Model struct {
	Substrings struct {
		Dlq          string   `yaml:"dlq"`
		Environments []string `yaml:"environments"`
	}
	Inputs []struct {
		Name   string      `yaml:"name"`
		Async  bool        `yaml:"async"`
		Values interface{} `yaml:"values"`
	}
	Enrichers []models.EnricherConfig
	Filters   []models.FilterConfig
	Outputs   []struct {
		Name string `yaml:"name"`
	}
}

func NewConfig() *Config {

	conf := Model{}

	if err := yaml.Unmarshal([]byte(config), &conf); err != nil {
		fmt.Printf("error loading config file config.yaml: %w", err)
		return nil
	}

	return &Config{Config: &conf}

}

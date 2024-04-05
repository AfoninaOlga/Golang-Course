package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Url string `yaml:"source_url"`
	DB  string `yaml:"db_file"`
}

func GetConfig(path string) (c Config, err error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &c)
	return
}

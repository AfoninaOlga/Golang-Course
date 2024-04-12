package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Url            string `yaml:"source_url"`
	DB             string `yaml:"db_file"`
	GoroutineCount int    `yaml:"parallel"`
}

func GetConfig(path string) (c Config, err error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &c)
	return
}

func ParseFlag() (configPath string) {
	flag.StringVar(&configPath, "c", "config.yaml", "flag sets config file path")
	flag.Parse()
	return
}

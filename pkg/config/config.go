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

func ParseFlag() (configPath string, sQuery string, useIndex bool) {
	flag.StringVar(&configPath, "c", "config.yaml", "flag sets config file path")
	flag.StringVar(&sQuery, "s", "", "flag sets search query")
	flag.BoolVar(&useIndex, "i", false, "flag sets index usage")
	flag.Parse()
	return
}

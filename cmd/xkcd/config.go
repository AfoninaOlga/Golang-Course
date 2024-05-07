package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Url            string `yaml:"source_url"`
	DB             string `yaml:"db_file"`
	GoroutineCount int    `yaml:"parallel"`
	Port           int    `yaml:"port"`
	Time           string `yaml:"update_time"`
}

func GetConfig(path string) (c Config, err error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &c)
	return
}

func ParseFlag() (configPath string, port int) {
	flag.StringVar(&configPath, "c", "config.yaml", "flag sets config file path")
	flag.IntVar(&port, "p", -1, "flag sets port for the server")
	flag.Parse()
	return
}

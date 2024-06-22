package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Api           string `yaml:"api_url"`
	Port          int    `yaml:"web_port"`
	TokenDuration uint   `yaml:"token_duration"`
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

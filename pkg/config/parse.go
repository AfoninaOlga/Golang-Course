package config

import (
	"flag"
	"math"
)

func ParseInput() (cnt uint, output bool, configPath string) {
	flag.UintVar(&cnt, "n", math.MaxUint, "flag takes number of comics to read")
	flag.BoolVar(&output, "o", false, "flag determines whether to print output or not")
	flag.StringVar(&configPath, "c", "config.yaml", "flag sets config file path")
	flag.Parse()
	return
}

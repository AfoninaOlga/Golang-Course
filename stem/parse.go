package main

import (
	"flag"
)

func ParseInput() (str string) {
	flag.StringVar(&str, "s", "", "flag takes string for stemming")
	flag.Parse()
	return
}

package utils

import (
	"flag"
)

func ParseInput() (cntIsSet bool, cnt uint, output bool) {
	flag.UintVar(&cnt, "n", 0, "flag takes number of comics to read")
	flag.BoolVar(&output, "o", false, "flag determines whether to print output or not")
	flag.Parse()
	flag.Visit(func(fl *flag.Flag) {
		if fl.Name == "n" {
			cntIsSet = true
		}
	})
	return
}

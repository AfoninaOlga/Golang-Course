package utils

import (
	"flag"
	"math"
)

func ParseInput() (cnt uint, output bool) {
	flag.UintVar(&cnt, "n", math.MaxUint, "flag takes number of comics to read")
	flag.BoolVar(&output, "o", false, "flag determines whether to print output or not")
	flag.Parse()
	return
}

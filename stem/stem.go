package main

import (
	"github.com/kljensen/snowball/english"
	"regexp"
	"strings"
)

func trimWord(word string) string {
	re := regexp.MustCompile(`[^a-zA-Z]`)
	return re.ReplaceAllString(word, "")
}

func StemInput(input string) string {
	var result []string
	stemmedWords := make(map[string]bool)
	re := regexp.MustCompile(`[-'.,!?:;]`)
	for _, s := range strings.Fields(re.ReplaceAllString(input, " ")) {
		s := english.Stem(trimWord(s), false)
		if len(s) <= 2 || english.IsStopWord(s) || stemmedWords[s] {
			continue
		}
		stemmedWords[s] = true
		result = append(result, s)
	}
	return strings.Join(result, " ")
}

package main

import (
	"github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

func StemInput(input string) string {
	var result []string
	stemmedWords := make(map[string]bool)
	f := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	for _, s := range strings.FieldsFunc(input, f) {
		s = english.Stem(s, false)
		if len(s) <= 2 || english.IsStopWord(s) || stemmedWords[s] {
			continue
		}
		stemmedWords[s] = true
		result = append(result, s)
	}
	return strings.Join(result, " ")
}

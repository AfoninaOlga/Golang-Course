package words

import (
	"fmt"
	"github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

func checkWord(word string) (ok bool, err error) {
	ok = true
	err = nil

	for _, c := range word {
		if !unicode.Is(unicode.Latin, c) {
			ok = false
			err = fmt.Errorf("unknown letter: %c", c)
		}
	}
	return
}

func StemInput(input string) ([]string, error) {
	var result []string
	input = strings.ReplaceAll(input, `\n`, " ")
	stemmedWords := make(map[string]bool)
	f := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	for _, s := range strings.FieldsFunc(input, f) {
		ok, err := checkWord(s)
		if !ok {
			return result, err
		}
		s = english.Stem(s, false)
		if len(s) <= 2 || english.IsStopWord(s) || s == "alt" || stemmedWords[s] {
			continue
		}
		stemmedWords[s] = true
		result = append(result, s)
	}
	return result, nil
}

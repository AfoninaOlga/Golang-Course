package stemmer

import (
	"errors"
	"fmt"
	"github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

func checkWord(word string) error {
	for _, c := range word {
		if !unicode.Is(unicode.Latin, c) {
			return fmt.Errorf("unknown letter: %c ", c)
		}
	}
	return nil
}

func Stem(input string) ([]string, error) {
	var result []string
	var err error
	stemmedWords := make(map[string]bool)
	f := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	for _, s := range strings.FieldsFunc(input, f) {
		if checkErr := checkWord(s); checkErr != nil {
			err = errors.Join(err, checkErr)
			continue
		}
		s = english.Stem(s, false)
		// added "alt" to stop-list cause transcript may contain "Alt:<alternative description>"
		if len(s) <= 2 || english.IsStopWord(s) || s == "alt" || stemmedWords[s] {
			continue
		}
		stemmedWords[s] = true
		result = append(result, s)
	}
	return result, err
}

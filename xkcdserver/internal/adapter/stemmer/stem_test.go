package stemmer

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckWord(t *testing.T) {
	testTable := []struct {
		word     string
		expected error
	}{
		{
			word:     "apple",
			expected: nil,
		},
		{
			word:     "aaaβήτα",
			expected: fmt.Errorf("unknown letter: β "),
		},
	}
	for _, testCase := range testTable {
		err := checkWord(testCase.word)
		assert.Equal(t, err, testCase.expected)
	}
}

func TestStem(t *testing.T) {
	testTable := []struct {
		text        string
		expectedRes []string
		expectedErr error
	}{
		{
			text:        "apple, a~ day",
			expectedRes: []string{"appl", "day"},
			expectedErr: nil,
		},
		{
			text:        "apple aaaβήτα day",
			expectedRes: []string{"appl", "day"},
			expectedErr: errors.Join(nil, fmt.Errorf("unknown letter: β ")),
		},
		{
			text:        "",
			expectedRes: nil,
			expectedErr: nil,
		},
	}
	for _, testCase := range testTable {
		res, err := Stem(testCase.text)
		assert.Equal(t, res, testCase.expectedRes)
		assert.Equal(t, err, testCase.expectedErr)
	}
}

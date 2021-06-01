package trigram

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractTrigrams(t *testing.T) {

	var tests = []struct {
		input    string
		trigrams []string
	}{
		{
			input:    "word",
			trigrams: []string{" w", " wo", "wor", "ord", "rd "},
		},
		{
			input:    "two words",
			trigrams: []string{" t", " tw", "two", "wo ", " w", " wo", "wor", "ord", "rds", "ds "},
		},
		{
			input:    "a",
			trigrams: []string{" a", " a "},
		},
		{
			input:    "    a       ",
			trigrams: []string{" a", " a "},
		},
		{
			input: "",
		},
		{
			input: "          ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {

			trigrams := ExtractTrigrams(tt.input)
			require.Equal(t, len(tt.trigrams), len(trigrams))

			for i, tg := range trigrams {
				assert.Equal(t, tt.trigrams[i], fmt.Sprint(tg))
			}
		})
	}
}

func TestStringsSimilarity(t *testing.T) {

	var tests = []struct {
		s1       string
		s2       string
		expected float64
	}{
		{
			s1:       "word",
			s2:       "two words",
			expected: 0.363636,
		},
		{
			s1:       "1600 Pennsylvania Ave",
			s2:       "1600 Penna Avenue",
			expected: 0.428571,
		},
	}

	for _, tt := range tests {
		t.Run(tt.s1, func(t *testing.T) {

			similarity := StringsSimilarity(tt.s1, tt.s2)
			assert.InDelta(t, tt.expected, similarity, .0001)
		})
	}
}

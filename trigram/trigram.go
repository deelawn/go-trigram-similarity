// trigram is meant to achieve feature parity with postgresql's similarity function. It works best with ascii
// characters; implementing trigram on a rune basis is outside the scope of our current needs since all of
// our customers store there address in a Latin script.
package trigram

import "strings"

const (
	spaceChar   byte  = 32
	notYetFound uint8 = 0
	foundInT1   uint8 = 1
	foundInBoth uint8 = 2
)

// TrigramsSimilarity returns how similar two slices of Trigrams are.
func TrigramsSimilarity(t1, t2 Trigrams) float64 {

	unique := map[Trigram]uint8{}

	// Do this first so we are easily able to retrieve the number of unique trigrams.
	for _, t := range append(t1, t2...) {
		unique[t] = notYetFound
	}

	for _, t := range t1 {
		if unique[t] == notYetFound {
			unique[t]++
		}
	}

	for _, t := range t2 {
		if unique[t] == foundInT1 {
			unique[t]++
		}
	}

	numUnique := len(unique)
	var numMatched int

	for _, found := range unique {
		if found == foundInBoth {
			numMatched++
		}
	}

	return float64(numMatched) / float64(numUnique)
}

// StringsSimilarity returns how similar two strings are using Trigrams.
func StringsSimilarity(s1, s2 string) float64 {
	return TrigramsSimilarity(ExtractTrigrams(s1), ExtractTrigrams(s2))
}

// Trigram (single) three uint8 --> uint32 (shifted)
type Trigram uint32

// String is great to use in the tests to confirm things were tokenized correctly while
// keeping them readable.
func (t Trigram) String() string {

	var b []byte
	if byte(t) == 0 { // Account for the case where we only have two chars.
		b = []byte{byte(t >> 16), byte(t >> 8)}
	} else {
		b = []byte{byte(t >> 16), byte(t >> 8), byte(t)}
	}

	return string(b[:])
}

// Trigrams is a slice of Trigrams
type Trigrams []Trigram

func (t Trigrams) TrigramsSimilarity(trigrams Trigrams) float64 {
	return TrigramsSimilarity(t, trigrams)
}

func (t Trigrams) StringSimilarity(s string) float64 {

	trigrams := ExtractTrigrams(s)
	return TrigramsSimilarity(t, trigrams)
}

// ExtractTrigrams extracts and returns all trigrams from a string
func ExtractTrigrams(s string) Trigrams {

	var trigrams Trigrams

	// Split by whitespace and construct the trigrams for each token.
	for _, s := range strings.Fields(s) {
		trigrams = append(trigrams, extractTrigrams(s)...)
	}
	return trigrams
}

// extractTrigrams pads the input with spaces and extracts and returns trigrams
func extractTrigrams(s string) Trigrams {

	var trigrams Trigrams

	padded := []byte{spaceChar, spaceChar} // pad it
	padded = append(padded, []byte(s)...)
	padded = append(padded, spaceChar) // pad it

	for i := 0; i < len(padded)-2; i++ {

		var trigram Trigram

		if i == 0 { // We only want, at most, one leading space
			trigram = Trigram(uint32(padded[i+1])<<16 | uint32(padded[i+2])<<8)
		} else {
			trigram = Trigram(uint32(padded[i])<<16 | uint32(padded[i+1])<<8 | uint32(padded[i+2]))
		}

		trigrams = append(trigrams, trigram)
	}

	return trigrams
}

package main

import (
	"reflect"
	"testing"
)

func newTestDict() Dictionary {
	return NewDictionaryWithName(jmdictSnippet)
}

func TestParseSentence(t *testing.T) {
	testCases := []struct {
		input             string
		expEntries        int
		expEntriesWithDef int
	}{
		{"", 0, 0},
		{"。", 1, 0},
		{"latin", 5, 0},
		{"そうか！なるほどね", 4, 3},
		{"パンを食べる", 3, 3},
		{"パンをたべた", 3, 3},
	}

	dict := newTestDict()
	for _, tc := range testCases {
		res := dict.ParseSentence(tc.input)
		len := len(res)
		if len != tc.expEntries {
			t.Errorf("Should return %d entries for input %q, not %d:\n%s", tc.expEntries, tc.input, len, res)
		}
	}
}

func TestLookupConjugations(t *testing.T) {
	testCases := []struct {
		input      string
		expMatches int
	}{
		{
			input:      "食べた",
			expMatches: 1,
		},
		{
			input:      "走た", // godan shouldn't match
			expMatches: 0,
		},
	}

	dict := newTestDict()
	// overwrite conjugations so as to not be dependent on the contents of the actual file
	dict.conjugations = []Conjugation{
		{
			Ending: "た",
			Base:   "る",
			POS:    "Ichidan verb",
			Name:   "Past form",
		},
	}

	for _, tc := range testCases {
		res := dict.lookupConjugations(tc.input)
		if len(res) != tc.expMatches {
			t.Errorf("Expected %d matches but got %d", tc.expMatches, len(res))
		}
	}
}

func TestNextDelimiter(t *testing.T) {
	testCases := []struct {
		input     string
		fromIndex int
		exp       int
	}{
		{"", 0, 0},
		{"。", 0, 0},
		{"this is some text", 0, 17},
		{"this is some text。abc defg", 0, 17},
		{"ああ、そうか！なるほどね", 0, 6 * 3},
		{"東京で住んでいる", 0, 8 * 3},
		{"東京で住んでいますか？それとも京都で？", 0, 10 * 3},
		{"東京で住んでいますか？それとも京都で？", 11 * 3, 18 * 3},
	}

	for _, tc := range testCases {
		res := nextDelimiter(tc.input, tc.fromIndex)
		if res != tc.exp {
			t.Errorf("Should have returned %v, not %v for input %q with fromIndex %d",
				tc.exp, res, tc.input, tc.fromIndex)
		}
	}
}

func TestNextNonDelimiter(t *testing.T) {
	testCases := []struct {
		input     string
		fromIndex int
		expIndex  int
		expWord   *Word
	}{
		{"", 0, 0, nil},
		{"。", 0, 3, &Word{Original: "。"}},
		{"this is some text", 0, 0, nil},
		{"ああ、そうか！なるほどね", 0, 0, nil},
		{"」？！東京で住んでいる", 0, 3 * 3, &Word{Original: "」？！"}},
	}

	for _, tc := range testCases {
		resWord, resIndex := nextNonDelimiter(tc.input, tc.fromIndex)
		if resIndex != tc.expIndex {
			t.Errorf("Should have returned index %v, not %v for input %q with fromIndex %d",
				tc.expIndex, resIndex, tc.input, tc.fromIndex)
		}
		if !reflect.DeepEqual(resWord, tc.expWord) {
			t.Errorf("Should have returned word %+v, not %+v for input %q with fromIndex %d",
				tc.expWord, resWord, tc.input, tc.fromIndex)
		}
	}
}

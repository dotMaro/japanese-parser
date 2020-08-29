package main

import "testing"

func TestParseSentence(t *testing.T) {
	testCases := []struct {
		input             string
		expEntries        int
		expEntriesWithDef int
	}{
		{"", 0, 0},
		{"。", 1, 0},
		{"latin", 5, 0},
		{"ああ、そうか！なるほどね", 5, 4},
		{"パンを食べる", 3, 3},
		{"パンをたべる", 3, 3},
		// {"東京で住んでいる", 0, 8 * 3}, TODO: Enable these when conjugation is supported
		// {"東京で住んでいますか？それとも京都で？", 0, 10 * 3},
		// {"東京で住んでいますか？それとも京都で？", 11 * 3, 18 * 3},
	}

	dict := NewDictionary()
	for _, tc := range testCases {
		res := dict.ParseSentence(tc.input)
		len := len(res)
		if len != tc.expEntries {
			t.Errorf("Should return %d entries for input %q, not %d:\n%s", tc.expEntries, tc.input, len, res)
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

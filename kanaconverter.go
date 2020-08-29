package main

import "strings"

// KataToHira converts katakana to hiragana.
func KataToHira(i string) string {
	var converted strings.Builder
	for _, r := range i {
		converted.WriteRune(kataToHiraRune(r))
	}
	return converted.String()
}

// kataToHiraRune returns the hiragana equivalent of the katakana rune.
// If the input is not katakana then it is returned unchanged.
func kataToHiraRune(r rune) rune {
	if isKata(r) {
		return r - 0x60
	}
	return r
}

func isKata(r rune) bool {
	return r >= 'ァ' && r <= 'ヶ'
}

// hiraToKataRune returns the katakana equivalent of the hiragana rune.
// If the input is not hiragana then it is returned unchanged.
func hiraToKataRune(r rune) rune {
	if isHira(r) {
		return r + 0x60
	}
	return r
}

func isHira(r rune) bool {
	return r >= 'ぁ' && r <= 'ゖ'
}

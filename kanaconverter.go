package japaneseparser

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
func kataToHiraRune(i rune) rune {
	if i >= 'ァ' && i <= 'ヶ' {
		return i - 0x60
	}
	return i
}

// hiraToKataRune returns the katakana equivalent of the hiragana rune.
// If the input is not hiragana then it is returned unchanged.
func hiraToKataRune(i rune) rune {
	if i >= 'ぁ' && i <= 'ゖ' {
		return i + 0x60
	}
	return i
}

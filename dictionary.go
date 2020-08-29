package main

import (
	"strings"
	"unicode/utf8"
)

// Dictionary is used to look up words.
type Dictionary struct {
	entries map[string]Entry
}

// NewDictionary reads and parses JMDict and returns a Dictionary.
func NewDictionary() Dictionary {
	jmdict := ReadJMDict()

	entries := make(map[string]Entry)
	for _, entry := range jmdict.Entries {
		// TODO: Make into a slice and put particles at top
		for _, kanji := range entry.Kanji {
			entries[kanji] = entry
		}

		for _, kana := range entry.Readings {
			entries[kana] = entry
		}
	}

	return Dictionary{entries: entries}
}

// LookupWord word s and return an entry, or nil if no match is found.
func (d *Dictionary) LookupWord(s string) Word {
	entry, ok := d.entries[s]
	if ok {
		return Word{
			Original:   s,
			Definition: &entry,
		}
	}
	return Word{Original: s}
}

// Sentence is a slice of Words.
type Sentence []Word

func (s Sentence) String() string {
	var b strings.Builder
	for i, w := range s {
		b.WriteString(w.String())
		if i != len(s)-1 {
			b.WriteRune('\n')
		}
	}
	return b.String()
}

// Word contains the original text and its definition.
type Word struct {
	Original   string
	Definition *Entry
}

func (w Word) String() string {
	if w.Definition != nil {
		return w.Definition.String()
	}
	return w.Original
}

// ParseSentence and returns a Sentence,
func (d *Dictionary) ParseSentence(s string) Sentence {
	var words Sentence
	punctuationWord, start := firstNonDelimiter(s, 0)
	if punctuationWord != nil {
		words = []Word{*punctuationWord}
	} else {
		words = make([]Word, 0)
	}
	nextDelim := nextDelimiter(s, start)
	end := nextDelim
	for start != -1 && start < len(s) {
		lookup := s[start:end]
		word := d.LookupWord(lookup)
		if word.Definition == nil {
			_, endRuneLen := utf8.DecodeLastRuneInString(lookup)
			if end-endRuneLen > start {
				end -= endRuneLen
			} else {
				// No words found with current start, give up on it
				words = append(words, word)
				_, runeLen := utf8.DecodeRuneInString(lookup)
				start += runeLen
				end = nextDelim
			}
		} else {
			var punctuationWord *Word
			punctuationWord, start = firstNonDelimiter(s, end)
			if punctuationWord != nil {
				words = append(words, word, *punctuationWord)
			} else {
				words = append(words, word)
			}
			nextDelim = nextDelimiter(s, start+1)
			end = nextDelim
		}
	}

	return words
}

type window struct {
	start     int
	end       int
	nextDelim int
	fullText  string
}

func (w *window) curText() string {
	return w.fullText[w.start:w.end]
}

func (w *window) narrow(s string) {
	if w.end > w.start {
		_, runeLen := utf8.DecodeLastRuneInString(s)
		if w.end-runeLen <= w.start {

		}
		w.end -= runeLen
	} else {
		// No words found with current start, give up on it
		// words = append(words, word)
		_, runeLen := utf8.DecodeRuneInString(s)
		w.start += runeLen
		w.end = w.nextDelim
	}
}

func (w *window) incrementStart() {
	_, runeLen := utf8.DecodeRuneInString(w.curText())
	w.start += runeLen
	w.end = w.nextDelim
}

const delimiters = "。！？「」（）"

func isDelimiter(r rune) bool {
	for _, d := range delimiters {
		if r == d {
			return true
		}
	}
	return false
}

// nextDelimiter returns the index of the next delimiter,
// or the last position if none was found.
func nextDelimiter(s string, fromIndex int) int {
	len := len(s)
	if fromIndex >= len {
		return len
	}
	for i, c := range s[fromIndex:] {
		if isDelimiter(c) {
			return i + fromIndex
		}
	}
	return len
}

// firstNonDelimiter returns the index of the first byte that's
// not a delimiter. If there are no non-delimiters in the string then -1 is returned.
func firstNonDelimiter(s string, fromIndex int) (*Word, int) {
	if len(s) == fromIndex {
		return nil, len(s)
	}

	for i, c := range s[fromIndex:] {
		if !isDelimiter(c) {
			var word *Word
			if i > fromIndex {
				word = &Word{Original: s[fromIndex:i]}
			}
			return word, i + fromIndex
		}
	}

	word := Word{Original: s[fromIndex:]}
	return &word, len(s)
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Dictionary is used to look up words.
type Dictionary struct {
	entries      map[string][]JMdictEntry
	conjugations []Conjugation
}

// Conjugation contains information about a word inflection.
type Conjugation struct {
	Ending string `json:"ending"`
	Base   string `json:"base"`
	POS    string `json:"pos"` // part of speech
	Name   string `json:"name"`
}

// NewDictionary reads and parses JMDict and returns a Dictionary.
func NewDictionary() Dictionary {
	jmdict := ReadJMDict()

	entries := make(map[string][]JMdictEntry)
	for _, entry := range jmdict.Entries {
		hasKanjiReadings := len(entry.Kanji) > 0
		for _, kana := range entry.Readings {
			defs, ok := entries[kana]
			if !ok {
				entries[kana] = []JMdictEntry{entry}
			} else {
				// Prepend if it has no kanji readings
				if !hasKanjiReadings {
					entries[kana] = append([]JMdictEntry{entry}, defs...)
				} else {
					entries[kana] = append(defs, entry)
				}
			}
		}

		for _, kanji := range entry.Kanji {
			defs, ok := entries[kanji]
			if !ok {
				entries[kanji] = []JMdictEntry{entry}
			} else {
				entries[kanji] = append(defs, entry)
			}
		}
	}

	return Dictionary{
		entries:      entries,
		conjugations: readConjugations(jmdict.Entities),
	}
}

func readConjugations(entities map[string]string) []Conjugation {
	file, err := os.Open("./conjugations.json")
	if err != nil {
		panic(fmt.Sprintf("failed to open conjugations file: %v", err))
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	var conjugations []Conjugation
	err = dec.Decode(&conjugations)
	if err != nil {
		panic(fmt.Sprintf("failed to parse conjugations file: %v", err))
	}

	for i := 0; i < len(conjugations); i++ {
		c := conjugations[i]
		pos, ok := entities[c.POS]
		if !ok {
			panic(fmt.Sprintf("unknown pos %q for conjugation %+v", c.POS, c))
		}
		conjugations[i].POS = pos
	}

	return conjugations
}

// LookupWord word s and return an entry, or nil if no match is found.
func (d *Dictionary) LookupWord(s string) Word {
	word := Word{
		Original: s,
	}
	jmDictEntries, ok := d.entries[s]
	if ok {
		entries := make([]Entry, len(jmDictEntries))
		for i, e := range jmDictEntries {
			entries[i] = Entry{JMdictEntry: e}
		}

		word.Definitions = entries
	}

	// Check for conjugations
	conjugationEntries := d.lookupConjugations(s)
	if len(conjugationEntries) > 0 {
		if word.Definitions != nil {
			word.Definitions = append(word.Definitions, conjugationEntries...)
		} else {
			word.Definitions = conjugationEntries
		}
	}

	return word
}

// lookupConjugations for a string.
// Never returns a nil slice. If no entries were found it will be empty.
func (d *Dictionary) lookupConjugations(s string) []Entry {
	var entries []Entry

	for _, c := range d.conjugations {
		if strings.HasSuffix(s, c.Ending) {
			conjugationEntries, ok := d.entries[s[:len(s)-len(c.Ending)]+c.Base]
			if ok {
				for _, entry := range conjugationEntries {
					for _, sense := range entry.Sense {
						for _, pos := range sense.POS {
							if pos == c.POS {
								newEntries := make([]Entry, len(conjugationEntries))
								for i, e := range conjugationEntries {
									newEntries[i] = Entry{JMdictEntry: e, Conjugation: c}
								}
								if entries != nil {
									entries = append(entries, newEntries...)
								} else {
									entries = newEntries
								}
							}
						}
					}
				}
			}
		}
	}

	return entries
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
	Original    string
	Definitions []Entry
}

// Entry wraps the JMdict entry with additional information.
type Entry struct {
	JMdictEntry JMdictEntry
	Conjugation Conjugation
}

func (e Entry) String() string {
	return e.JMdictEntry.String() + " " + e.Conjugation.Name
}

func (w Word) String() string {
	if w.Definitions != nil && len(w.Definitions) > 0 {
		return w.Definitions[0].String()
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
		if word.Definitions == nil {
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

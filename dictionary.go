package main

// Dictionary is used to look up words.
type Dictionary struct {
	words map[string]Entry
}

// ParseDictionary reads and parses JMDict and returns a Dictionary.
func ParseDictionary() Dictionary {
	jmdict := ReadJMDict()

	words := make(map[string]Entry)
	for _, entry := range jmdict.Entries {
		for _, kanji := range entry.Kanji {
			words[kanji] = entry
		}

		for _, kana := range entry.Readings {
			words[kana] = entry
		}
	}

	return Dictionary{words: words}
}

// Lookup word s and return an entry, or nil if no match is found.
func (d *Dictionary) Lookup(s string) *Entry {
	entry, ok := d.words[s]
	if ok {
		return &entry
	}

	return nil
}

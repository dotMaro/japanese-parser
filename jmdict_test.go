package main

import "testing"

const jmdictSnippet = "JMdict_e-snippet"

func TestReadJMDictWithName(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	// Use the JMdict snippet since it is checked into the repo,
	// as opposed to the full dictionary.
	dict := ReadJMDictWithName(jmdictSnippet)

	len := len(dict.Entries)
	const minExpectedEntries = 7
	if len < minExpectedEntries {
		t.Errorf("Should have at least %d entries but had %d", minExpectedEntries, len)
	}
}

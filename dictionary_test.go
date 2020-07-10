package japaneseparser

import "testing"

func TestReadJMDict(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	dict := ReadJMDict()

	len := len(dict.Entries)
	t.Logf("Len: %d\n", len)
	const minExpectedEntries = 180000
	if len < minExpectedEntries {
		t.Errorf("Should have at least %d entries but had %d", minExpectedEntries, len)
	}

	t.Log(dict.Entries[100000])

	t.Fail()
}

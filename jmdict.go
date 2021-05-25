package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// JMdict dictionary.
type JMdict struct {
	XMLName  xml.Name      `xml:"JMdict"`
	Entries  []JMdictEntry `xml:"entry"`
	Entities map[string]string
}

// JMdictEntry is an entry in JMdict.
type JMdictEntry struct {
	XMLName  xml.Name `xml:"entry" json:"-"`
	Kanji    []string `xml:"k_ele>keb" json:"kanji"`
	Readings []string `xml:"r_ele>reb" json:"readings"`
	Sense    []Sense  `xml:"sense" json:"sense"`
}

// Sense in JMdict.
type Sense struct {
	Glossary []string `xml:"gloss" json:"glossary"`
	POS      []string `xml:"pos" json:"pos"` // part of speech
}

func (s Sense) String() string {
	var b strings.Builder
	for i, g := range s.Glossary {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(g)
	}
	return b.String()
}

func (e JMdictEntry) String() string {
	return fmt.Sprintf("%v %v %v", e.Kanji, e.Readings, e.Sense[0].Glossary[0])
}

// DetailedString returns a more detailed string with all meanings included.
func (e JMdictEntry) DetailedString() string {
	var b strings.Builder
	for i, kanji := range e.Kanji {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(kanji)
	}
	if len(e.Kanji) > 0 {
		b.WriteRune('\n')
	}
	for i, kana := range e.Readings {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(kana)
	}
	b.WriteRune('\n')
	for i, s := range e.Sense {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, s.String()))
	}
	return b.String()
}

// IsParticle returns true if the JMdictEntry has a particle part of speech.
func (e JMdictEntry) IsParticle() bool {
	for _, s := range e.Sense {
		for _, pos := range s.POS {
			if pos == "particle" {
				return true
			}
		}
	}
	return false
}

const defaultJMdictName = "JMdict_e"

// ReadJMDict reads and parses the default JMdict file.
// Upon error it will panic.
func ReadJMDict() JMdict {
	return ReadJMDictWithName(defaultJMdictName)
}

// ReadJMDictWithName reads and parses a specified JMdict file.
// Upon error it will panic.
func ReadJMDictWithName(name string) JMdict {
	file, err := os.Open(fmt.Sprintf("./dict/%s", name))
	if err != nil {
		panic(fmt.Sprintf("unable to read JMdict file: %v", err))
	}
	defer file.Close()

	var dict JMdict

	dec := xml.NewDecoder(file)
	entities, err := parseEntityCodes(file)
	if err != nil {
		panic(fmt.Sprintf("unable to parse entities in JMdict file: %v", err))
	}
	dec.Entity = entities
	dict.Entities = entities

	err = dec.Decode(&dict)
	if err != nil {
		panic(fmt.Sprintf("unable to parse JMdict file: %v", err))
	}

	return dict
}

// parseEntityCodes in DTD because apparently the standard decoder does not.
func parseEntityCodes(dict io.ReadSeeker) (map[string]string, error) {
	defer dict.Seek(0, io.SeekStart)

	entities := make(map[string]string)
	scan := bufio.NewReader(dict)
	for {
		entity, err := scan.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if entity == "<JMdict>\n" {
			// Once the dict entries actually start, all entities should
			// be listed already, so abort entity lookup.
			break
		}

		if strings.HasPrefix(entity, "<!ENTITY ") {
			splitEntity := strings.Split(entity, `"`)
			if len(splitEntity) != 3 {
				return nil, fmt.Errorf("unexpected entity format %q", entity)
			}
			// E.g. "<!ENTITY v5n " -> "v5n"
			entName := splitEntity[0][9 : len(splitEntity[0])-1]
			entDef := splitEntity[1]
			entities[entName] = entDef
		}
	}

	return entities, nil
}

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
	XMLName  xml.Name `xml:"entry"`
	Kanji    []string `xml:"k_ele>keb"`
	Readings []string `xml:"r_ele>reb"`
	Sense    []Sense  `xml:"sense"`
}

// Sense in JMdict.
type Sense struct {
	Glossary []string `xml:"gloss"`
	POS      []string `xml:"pos"` // part of speech
}

func (e JMdictEntry) String() string {
	return fmt.Sprintf("%v %v %v", e.Kanji, e.Readings, e.Sense[0].Glossary[0])
}

// ReadJMDict reads and parses a JMdict file.
// Upon error it will panic.
func ReadJMDict() JMdict {
	file, err := os.Open("./dict/JMdict_e")
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

package main

func main() {
	// dict := NewDictionaryWithName("JMdict_e-snippet")
	dict := NewDictionary()
	// fmt.Printf("%v\n", dict.ParseSentence("日本語上手"))
	// fmt.Printf("%v\n", dict.LookupWord("流れ").DetailedString())
	StartServer(dict)
}

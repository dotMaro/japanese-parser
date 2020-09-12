package main

import "fmt"

func main() {
	dict := NewDictionary()
	fmt.Printf("%v\n", dict.ParseSentence("パンをたべた"))
	// fmt.Printf("%v\n", dict.LookupWord("が"))
}

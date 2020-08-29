package main

import "fmt"

func main() {
	dict := NewDictionary()
	fmt.Printf("%v\n", dict.ParseSentence("業務委託契約書の添削の依頼？"))
}

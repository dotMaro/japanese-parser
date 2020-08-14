package main

import "fmt"

func main() {
	dict := ParseDictionary()
	fmt.Printf("%v\n", dict.Lookup("こんにちは"))
}

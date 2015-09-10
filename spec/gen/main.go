package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("go run gen.go (bnf | json | go)")
		return
	}
	switch os.Args[1] {
	case "bnf":
		fromHTMLToBNF()
	case "json":
		fromBNFToJSON()
	case "go":
		fromBNFToGo()
	}
}

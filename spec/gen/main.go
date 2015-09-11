package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("gen (bnf | bnfj | goj | go | gof)")
		fmt.Println("bnf: from HTML to BNF")
		fmt.Println("bnfj: from BNF to BNF JSON")
		fmt.Println("goj: from BNF to Go JSON")
		fmt.Println("go: from BNF to Go")
		fmt.Println("gof: from BNF to Go funcs")
		return
	}
	switch os.Args[1] {
	case "bnf":
		fromHTMLToBNF()
	case "bnfj":
		fromBNFToJSON()
	case "goj":
		fromBNFToGoJSON()
	case "go":
		fromBNFToGo()
	case "gof":
		fromBNFToGoFuncs()
	}
}

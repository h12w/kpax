package main

import (
	"fmt"
	"os"

	"h12.me/wipro"
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
	bnfFile := "bnf.txt"
	switch os.Args[1] {
	case "bnf":
		fromHTMLToBNF()
	case "bnfj":
		bnf := wipro.ParseBNF(bnfFile)
		fmt.Println(bnf.JSON())
	case "goj":
		bnf := wipro.ParseBNF(bnfFile)
		fmt.Println(bnf.GoTypes().JSON())
	case "go":
		bnf := wipro.ParseBNF(bnfFile)
		goTypes := bnf.GoTypes().RemoveDecl("RequestMessage")
		goTypes.Fprint(os.Stdout)
	case "gof":
		bnf := wipro.ParseBNF(bnfFile)
		goTypes := bnf.GoTypes()
		goTypes.GoFuncs(os.Stdout)
	}
}

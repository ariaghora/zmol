package main

import (
	"fmt"
	"os"

	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/parser"
)

// The interpreter
type Zmol struct {
}

type ZValue struct {
}

func (z *Zmol) Run(code string) *ZValue {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("unexpected problem encountered, aborting")
			os.Exit(1)
		}
	}()

	lexer := lexer.NewLexer(code)
	err := lexer.Lex()

	parser := parser.NewParser(lexer)
	program := parser.ParseProgram()

	fmt.Println(program.Str())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func printBanner() {
	fmt.Println("Zmol 0.0.1")
}

func main() {
	printBanner()
	for {
		var code string
		fmt.Print("> ")
		fmt.Scanln(&code)

		switch code {
		case "exit":
			os.Exit(0)
		default: // Run the code
			z := Zmol{}
			z.Run(code)
		}
	}
}

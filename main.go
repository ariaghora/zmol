package main

import (
	"fmt"
	"os"

	"github.com/ariaghora/zmol/zmol/lexer"
)

// The interpreter
type Zmol struct {
}

type ZValue struct {
}

func (z *Zmol) Run(code string) *ZValue {
	lexer := lexer.NewLexer(code)
	err := lexer.Lex()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func main() {}

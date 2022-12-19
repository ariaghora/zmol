package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/parser"
	"github.com/fatih/color"
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
	result := eval.Eval(program)

	fmt.Println(result.Str())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

var Banner = `
 ______     __    __     ______     __        
/\___  \   /\ "-./  \   /\  __ \   /\ \       
\/_/  /__  \ \ \-./\ \  \ \ \/\ \  \ \ \____  
  /\_____\  \ \_\ \ \_\  \ \_____\  \ \_____\ 
  \/_____/   \/_/  \/_/   \/_____/   \/_____/ 
`

func printBanner() {
	fmt.Println(Banner)
	fmt.Println("Zmol 0.0.1")
}

func main() {
	printBanner()
	for {
		color.Set(color.FgHiMagenta)
		fmt.Print(">>> ")
		color.Unset()

		in := bufio.NewReader(os.Stdin)
		code, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		code = strings.TrimSpace(code)

		switch code {
		case ".exit":
			os.Exit(0)
		default: // Run the code
			z := Zmol{}
			z.Run(code)
		}
	}
}

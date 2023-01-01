package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/ariaghora/zmol/pkg/compiler"
	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/native"
	"github.com/ariaghora/zmol/pkg/parser"
	"github.com/ariaghora/zmol/pkg/val"
	"github.com/ariaghora/zmol/pkg/vm"
	"github.com/fatih/color"
)

var Banner = `
 ______     __    __     ______     __        
/\___  \   /\ "-./  \   /\  __ \   /\ \       
\/_/  /__  \ \ \-./\ \  \ \ \/\ \  \ \ \____  
  /\_____\  \ \_\ \ \_\  \ \_____\  \ \_____\ 
  \/_____/   \/_/  \/_/   \/_____/   \/_____/ 
`

// The interpreter
type Zmol struct {
	vm    *vm.VM
	state *eval.ZmolState
}

func NewZmol() *Zmol {
	return &Zmol{
		state: eval.NewZmolState(nil),
	}
}

func (z *Zmol) Run(code string) (val.ZValue, error) {
	lexer := lexer.NewLexer(code)
	err := lexer.Lex()
	if err != nil {
		return nil, err
	}
	parser := parser.NewParser(lexer)
	program, err := parser.ParseProgram()
	if err != nil {
		return nil, err
	}

	compiler := compiler.NewCompiler()
	err = compiler.Compile(program)
	if err != nil {
		return nil, err
	}

	bytecode := compiler.Bytecode()
	z.vm = vm.NewVM(bytecode)
	err = z.vm.Run()

	if err != nil {
		return nil, err
	}

	result := z.vm.LastPoppedStackElem()
	return result, nil
}

//go:embed VERSION
var Version string

func printBanner() {
	fmt.Println(Banner)
	fmt.Println("Zmol ", Version)
}

func main() {
	z := NewZmol()

	nativeFuncRegistry := native.NewNativeFuncRegistry(z.state)
	nativeFuncRegistry.RegisterNativeFunc()

	if len(os.Args) > 1 {
		fileName := os.Args[1]
		sourceCode, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		z.state.Eval(string(sourceCode))
	} else {
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
				val, err := z.Run(code)
				if err != nil {
					fmt.Println("ERROR:", err)
					continue
				}

				if val == nil {
					fmt.Println("No value returned possibly because of interpreter error")
					continue
				}
			}
		}
	}
}

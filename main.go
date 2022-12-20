package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ariaghora/zmol/pkg/eval"
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
	state *eval.ZmolState
}

func NewZmol() *Zmol {
	return &Zmol{
		state: eval.NewZmolState(nil),
	}
}

func (z *Zmol) Run(code string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("unexpected problem encountered, aborting")
			os.Exit(1)
		}
	}()

	// state := eval.NewZmolState()
	result := z.state.Eval(code)

	fmt.Println(result.Str())
	// fmt.Println(result)

}

func printBanner() {
	fmt.Println(Banner)
	fmt.Println("Zmol 0.0.1")
}

func main() {
	printBanner()
	z := NewZmol()
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
			z.Run(code)
		}
	}
}

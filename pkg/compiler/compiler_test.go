package compiler

import (
	"fmt"
	"testing"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/bytecode"
	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/parser"
	"github.com/ariaghora/zmol/pkg/val"
)

func TestIntegerArith(t *testing.T) {
	tests := []struct {
		Input                string
		expectedConstants    []interface{}
		expectedInstructions []bytecode.Instructions
	}{
		{
			Input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
			},
		},
	}

	for _, test := range tests {
		program := parse(test.Input)
		compiler := NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()
		err = testInstructions(t, bytecode, test.expectedInstructions)
		if err != nil {
			t.Fatalf("testInstructions error: %s", err)
		}

		err = testConstants(t, bytecode, test.expectedConstants)
		if err != nil {
			t.Fatalf("testConstants error: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	err := l.Lex()
	if err != nil {
		panic(err)
	}
	p := parser.NewParser(l)
	res, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}

	return res
}

func testInstructions(t *testing.T, bytecode *Bytecode, expected []bytecode.Instructions) error {
	concatenated := concatInstructions(expected)

	if len(bytecode.Instructions) != len(concatenated) {
		t.Fatalf("wrong instructions length. want=%q, got=%q", concatenated, bytecode.Instructions)
	}

	for i, ins := range concatenated {
		if bytecode.Instructions[i] != ins {
			t.Fatalf("wrong instruction at %d. want=%d, got=%d", i, ins, bytecode.Instructions[i])
		}
	}

	return nil
}

func concatInstructions(inss []bytecode.Instructions) bytecode.Instructions {
	out := bytecode.Instructions{}
	for _, ins := range inss {
		out = append(out, ins...)
	}

	return out
}

func testConstants(t *testing.T, bytecode *Bytecode, expected []interface{}) error {
	constants := bytecode.Constants
	if len(constants) != len(expected) {
		t.Fatalf("wrong constants length. want=%d, got=%d", len(expected), len(constants))
	}

	for i, exp := range expected {
		switch exp := exp.(type) {
		case int:
			err := testIntegerObject(t, constants[i], int64(exp))
			if err != nil {
				return fmt.Errorf("testIntegerObject error: %s", err)
			}
		}
	}

	return nil
}

func testIntegerObject(t *testing.T, obj val.ZValue, expected int64) error {
	result, ok := obj.(*val.ZInt)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", obj, obj)
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value. want=%d, got=%d", expected, result.Value)
	}

	return nil
}

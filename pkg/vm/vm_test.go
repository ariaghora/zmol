package vm

import (
	"fmt"
	"testing"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/compiler"
	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/parser"
	"github.com/ariaghora/zmol/pkg/val"
)

func TestIntegerArith(t *testing.T) {
	tests := []VMTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"2 * 2", 4},
		{"4 / 2", 2},
		{"5 % 2", 1},
		{"4 % 2", 0},
	}

	runVMTests(t, tests)
}

func TestFloatArith(t *testing.T) {
	tests := []VMTestCase{
		{"1.0", 1.0},
		{"2.0", 2.0},
		{"1.0 + 2.0", 3.0},
		{"1.0 - 2.0", -1.0},
		{"2.0 * 2.0", 4.0},
		{"4.0 / 2.0", 2.0},
		{"5.5 % 2.0", 1.5},
	}

	runVMTests(t, tests)
}

func parse(source string) *ast.Program {
	l := lexer.NewLexer(source)
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

type VMTestCase struct {
	input    string
	expected interface{}
}

func runVMTests(t *testing.T, tests []VMTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.NewCompiler()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compile error: %s", err)
		}

		vm := NewVM(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("runtime error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedValue(t, stackElem, tt.expected)
	}
}

func testExpectedValue(t *testing.T, value val.ZValue, expected interface{}) {
	t.Helper()
	var err error
	switch expected := expected.(type) {
	case int:
		err = testIntegerValue(t, value, int64(expected))
	case float32, float64:
		err = testFloatValue(t, value, expected.(float64))
	}

	if err != nil {
		t.Fatal(err)
	}
}

func testIntegerValue(t *testing.T, value val.ZValue, expected int64) error {
	result, ok := value.(*val.ZInt)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", value, value)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testFloatValue(t *testing.T, value val.ZValue, expected float64) error {
	result, ok := value.(*val.ZFloat)
	if !ok {
		return fmt.Errorf("object is not Float. got=%T (%+v)", value, value)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
	}
	return nil
}

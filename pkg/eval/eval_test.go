package eval

import (
	"testing"

	"github.com/ariaghora/zmol/pkg/val"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1", 1},
		{"2", 2},
		{"123", 123},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionValue(t *testing.T) {
	input := "@(x){ x + 2 }"
	evaluated := testEval(input)
	fn, ok := evaluated.(*val.ZFunction)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Params)
	}

	if fn.Params[0].Str() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Params[0])
	}

	if fn.Body.Str() != "(x + 2)" {
		t.Fatalf("body is not (x + 2). got=%q", fn.Body.Str())
	}
}

func TestAnonFunctionCalls(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"@(x){ x + 2 }(2)", 4},
		{"@(x){ x + 2 }(3)", 5},
		{"@(x){ x + x }(4)", 8},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionCalls(t *testing.T) {
	source := `
	add_five = @(x){ x + 5 }
	add_five(5)
	`

	evaluated := testEval(source)
	testIntegerObject(t, evaluated, 10)

	sourceWithGlobal := `
	val = 5
	add_five = @(x){ x + val }
	add_five(5)
	`

	evaluated = testEval(sourceWithGlobal)
	testIntegerObject(t, evaluated, 10)
}

func testEval(input string) val.ZValue {
	state := NewZmolState(nil)
	value, err := state.Eval(input)
	if err != nil {
		return val.ERROR(err.Error())
	}
	return value
}

func testIntegerObject(t *testing.T, obj val.ZValue, expected int64) bool {
	result, ok := obj.(*val.ZInt)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

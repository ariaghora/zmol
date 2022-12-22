package lexer

import "testing"

func TestAssignmentStatement(t *testing.T) {
	input := `@x = 5.0
@y = 10
@foobar = 838383 `

	l := NewLexer(input)
	err := l.Lex()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tests := []ZTok{
		{Type: TokAt, Text: "@"},
		{Type: TokIdent, Text: "x"},
		{Type: TokAssign, Text: "="},
		{Type: TokFloat, Text: "5.0"},
		{Type: TokAt, Text: "@"},
		{Type: TokIdent, Text: "y"},
		{Type: TokAssign, Text: "="},
		{Type: TokInt, Text: "10"},
		{Type: TokAt, Text: "@"},
		{Type: TokIdent, Text: "foobar"},
		{Type: TokAssign, Text: "="},
		{Type: TokInt, Text: "838383"},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range l.Tokens {
		if tok != tests[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, tests[i], tok)
		}
	}
}

func TestFunctionDefinition(t *testing.T) {
	input := `@add(x, y) = x + y`

	l := NewLexer(input)
	err := l.Lex()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tests := []ZTok{
		{Type: TokAt, Text: "@"},
		{Type: TokIdent, Text: "add"},
		{Type: TokLParen, Text: "("},
		{Type: TokIdent, Text: "x"},
		{Type: TokComma, Text: ","},
		{Type: TokIdent, Text: "y"},
		{Type: TokRParen, Text: ")"},
		{Type: TokAssign, Text: "="},
		{Type: TokIdent, Text: "x"},
		{Type: TokPlus, Text: "+"},
		{Type: TokIdent, Text: "y"},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range l.Tokens {
		if tok != tests[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, tests[i], tok)
		}
	}
}

package lexer

import "testing"

func TestLexEmptySource(t *testing.T) {
	lexer := NewLexer("")
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.Tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(lexer.Tokens))
	}

	if lexer.Tokens[0].Type != TokEOF {
		t.Errorf("Expected EOF token, got %v", lexer.Tokens[0])
	}
}

func TestLexSimpleOpsSequence(t *testing.T) {
	lexer := NewLexer("+-><.")
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.Tokens) != 6 {
		t.Errorf("Expected 6 tokens, got %d", len(lexer.Tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokPlus, Text: "+"},
		{Type: TokMinus, Text: "-"},
		{Type: TokGt, Text: ">"},
		{Type: TokLt, Text: "<"},
		{Type: TokDot, Text: "."},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.Tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

func TestLexOpsSeqLeadingSpace(t *testing.T) {
	lexer := NewLexer(" +-><.")
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.Tokens) != 6 {
		t.Errorf("Expected 6 tokens, got %d", len(lexer.Tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokPlus, Text: "+"},
		{Type: TokMinus, Text: "-"},
		{Type: TokGt, Text: ">"},
		{Type: TokLt, Text: "<"},
		{Type: TokDot, Text: "."},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.Tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

func TestTwoCharOpsSeq(t *testing.T) {
	lexer := NewLexer("+ - <= >=")
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.Tokens) != 5 {
		t.Errorf("Expected 5 tokens, got %d", len(lexer.Tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokPlus, Text: "+"},
		{Type: TokMinus, Text: "-"},
		{Type: TokLTE, Text: "<="},
		{Type: TokGTE, Text: ">="},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.Tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

func TestInvalidTokenShouldFail(t *testing.T) {
	lexer := NewLexer("$")
	err := lexer.Lex()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestSimpleScript(t *testing.T) {
	source := `number1 = 5
	number2 = 10.0
	number3 = number1 + number2 `

	lexer := NewLexer(source)
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.Tokens) != 14 {
		t.Errorf("Expected 14 tokens, got %d", len(lexer.Tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokIdent, Text: "number1"},
		{Type: TokAssign, Text: "="},
		{Type: TokInt, Text: "5"},
		{Type: TokNewLine, Text: "\n"},
		{Type: TokIdent, Text: "number2"},
		{Type: TokAssign, Text: "="},
		{Type: TokFloat, Text: "10.0"},
		{Type: TokNewLine, Text: "\n"},
		{Type: TokIdent, Text: "number3"},
		{Type: TokAssign, Text: "="},
		{Type: TokIdent, Text: "number1"},
		{Type: TokPlus, Text: "+"},
		{Type: TokIdent, Text: "number2"},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.Tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

func TestTernarySymbols(t *testing.T) {
	source := `a ? b : c`

	lexer := NewLexer(source)
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.Tokens) != 6 {
		t.Errorf("Expected 6 tokens, got %d", len(lexer.Tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokIdent, Text: "a"},
		{Type: TokQuestion, Text: "?"},
		{Type: TokIdent, Text: "b"},
		{Type: TokLCurl, Text: ":"},
		{Type: TokIdent, Text: "c"},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.Tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

func TestCurlyBraces(t *testing.T) {
	source := `{ a = 5 }`

	lexer := NewLexer(source)
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.Tokens) != 6 {
		t.Errorf("Expected 6 tokens, got %d", len(lexer.Tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokLCurl, Text: "{"},
		{Type: TokIdent, Text: "a"},
		{Type: TokAssign, Text: "="},
		{Type: TokInt, Text: "5"},
		{Type: TokRCurl, Text: "}"},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.Tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

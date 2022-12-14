package lexer

import "testing"

func TestLexEmptySource(t *testing.T) {
	lexer := NewLexer("")
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(lexer.tokens))
	}

	if lexer.tokens[0].Type != TokEOF {
		t.Errorf("Expected EOF token, got %v", lexer.tokens[0])
	}
}

func TestLexSimpleOpsSequence(t *testing.T) {
	lexer := NewLexer("+-><.")
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.tokens) != 6 {
		t.Errorf("Expected 6 tokens, got %d", len(lexer.tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokPlus, Text: "+"},
		{Type: TokMinus, Text: "-"},
		{Type: TokGT, Text: ">"},
		{Type: TokLT, Text: "<"},
		{Type: TokDot, Text: "."},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.tokens {
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

	if len(lexer.tokens) != 6 {
		t.Errorf("Expected 6 tokens, got %d", len(lexer.tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokPlus, Text: "+"},
		{Type: TokMinus, Text: "-"},
		{Type: TokGT, Text: ">"},
		{Type: TokLT, Text: "<"},
		{Type: TokDot, Text: "."},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.tokens {
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

	if len(lexer.tokens) != 5 {
		t.Errorf("Expected 5 tokens, got %d", len(lexer.tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokPlus, Text: "+"},
		{Type: TokMinus, Text: "-"},
		{Type: TokLTE, Text: "<="},
		{Type: TokGTE, Text: ">="},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

func TestInvalidTokenShouldFail(t *testing.T) {
	lexer := NewLexer("!")
	err := lexer.Lex()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestSimpleScript(t *testing.T) {
	source := `number1 = 5
	number2 = 10.0
	number3 = number1 + number2
	`

	lexer := NewLexer(source)
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lexer.tokens) != 12 {
		t.Errorf("Expected 13 tokens, got %d", len(lexer.tokens))
	}

	expectedTokens := []ZTok{
		{Type: TokIdent, Text: "number1"},
		{Type: TokAssign, Text: "="},
		{Type: TokNumber, Text: "5"},
		{Type: TokIdent, Text: "number2"},
		{Type: TokAssign, Text: "="},
		{Type: TokNumber, Text: "10.0"},
		{Type: TokIdent, Text: "number3"},
		{Type: TokAssign, Text: "="},
		{Type: TokIdent, Text: "number1"},
		{Type: TokPlus, Text: "+"},
		{Type: TokIdent, Text: "number2"},
		{Type: TokEOF, Text: ""},
	}

	for i, tok := range lexer.tokens {
		if tok != expectedTokens[i] {
			t.Errorf("Expected token %d to be %v, got %v", i, expectedTokens[i], tok)
		}
	}
}

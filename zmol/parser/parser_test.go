package parser

import (
	"testing"

	"github.com/ariaghora/zmol/zmol/ast"
	"github.com/ariaghora/zmol/zmol/lexer"
)

func TestVarAssign(t *testing.T) {
	source := ` @a = 1
	@b = 2
	@c = 3 `

	lexer := lexer.NewLexer(source)
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
	}

	if program != nil {
		if len(program.Statements) != 3 {
			t.Errorf("Expected 3 statements, got %d", len(program.Statements))
		}
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"a"},
		{"b"},
		{"c"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if stmt.(*ast.VarrAssignmentStatement).Name.Value != tt.expectedIdentifier {
			t.Errorf("Expected identifier to be %s, got %s", tt.expectedIdentifier, stmt.(*ast.VarrAssignmentStatement).Name.Value)
		}
	}

}

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

func TestIdentifier(t *testing.T) {
	source := "test"
	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program := p.ParseProgram()

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.Identifier, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected expression to be *ast.Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "test" {
		t.Errorf("Expected identifier to be 'test', got %s", ident.Value)
	}
}

func TestIntegerLiteral(t *testing.T) {
	source := "5"
	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program := p.ParseProgram()

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected expression to be *ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("Expected literal to be 5, got %d", literal.Value)
	}
}

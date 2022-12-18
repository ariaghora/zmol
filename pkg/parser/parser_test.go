package parser

import (
	"testing"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/lexer"
)

func TestVarAssign(t *testing.T) {
	source := `let a = 1
	let b = 2
	let c = 3 `

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

func TestParseInfix(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		err := l.Lex()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		p := NewParser(l)
		program := p.ParseProgram()

		if program == nil {
			t.Errorf("ParseProgram() returned nil")
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected statement to be *ast.ExpressionStatement, got %T", program.Statements[0])
		}

		infix, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Errorf("Expected expression to be *ast.InfixExpression, got %T", stmt.Expression)
		}

		if infix.Operator != tt.operator {
			t.Errorf("Expected operator to be %s, got %s", tt.operator, infix.Operator)
		}

		left, ok := infix.Left.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("Expected left to be *ast.IntegerLiteral, got %T", infix.Left)
		}

		if left.Value != tt.left {
			t.Errorf("Expected left to be %d, got %d", tt.left, left.Value)
		}

		right, ok := infix.Right.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("Expected right to be *ast.IntegerLiteral, got %T", infix.Right)
		}

		if right.Value != tt.right {
			t.Errorf("Expected right to be %d, got %d", tt.right, right.Value)
		}

	}
}

func TestParseFuncLiteral(t *testing.T) {
	source := "@(x, y): x + y"
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

	function, ok := stmt.Expression.(*ast.FuncLiteral)
	if !ok {
		t.Errorf("Expected expression to be *ast.FuncLiteral, got %T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(function.Parameters))
	}

	if function.Parameters[0].Str() != "x" {
		t.Errorf("Expected parameter to be 'x', got %s", function.Parameters[0])
	}

	if function.Parameters[1].Str() != "y" {
		t.Errorf("Expected parameter to be 'y', got %s", function.Parameters[1])
	}

	if len(function.Body.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(function.Body.Statements))
	}

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.ExpressionStatement, got %T", function.Body.Statements[0])
	}

	infix, ok := body.Expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expected expression to be *ast.InfixExpression, got %T", body.Expression)
	}

	if infix.Operator != "+" {
		t.Errorf("Expected operator to be '+', got %s", infix.Operator)
	}

	left, ok := infix.Left.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected left to be *ast.Identifier, got %T", infix.Left)
	}

	if left.Value != "x" {
		t.Errorf("Expected left to be 'x', got %s", left.Value)
	}

	right, ok := infix.Right.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected right to be *ast.Identifier, got %T", infix.Right)
	}

	if right.Value != "y" {
		t.Errorf("Expected right to be 'y', got %s", right.Value)
	}

}

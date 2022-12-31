package parser

import (
	"testing"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/lexer"
)

func TestSimpleVarAssign(t *testing.T) {
	source := `a = 1
	`

	lexer := lexer.NewLexer(source)
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	parser := NewParser(lexer)
	program, _ := parser.ParseProgram()

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
	}

	if program != nil {
		if len(program.Statements) != 1 {
			t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		}
	}

}

func TestVarAssign(t *testing.T) {
	source := `a = 1
	b = 2
	c = 3 `

	lexer := lexer.NewLexer(source)
	err := lexer.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	parser := NewParser(lexer)
	program, _ := parser.ParseProgram()

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
		// expression statement with = operator
		stmt := program.Statements[i]

		if stmt.(*ast.ExpressionStatement).Expression.(*ast.InfixExpression).Operator != "=" {
			t.Errorf("Expected operator to be =, got %s", stmt.(*ast.ExpressionStatement).Expression.(*ast.InfixExpression).Operator)
		}

		// left side of expression statement
		left := stmt.(*ast.ExpressionStatement).Expression.(*ast.InfixExpression).Left
		if left.(*ast.Identifier).Value != tt.expectedIdentifier {
			t.Errorf("Expected identifier to be %s, got %s", tt.expectedIdentifier, left.(*ast.Identifier).Value)
		}

		// right side of expression statement
		right := stmt.(*ast.ExpressionStatement).Expression.(*ast.InfixExpression).Right
		if right.(*ast.IntegerLiteral).Value != int64(i+1) {
			t.Errorf("Expected integer to be %d, got %d", i+1, right.(*ast.IntegerLiteral).Value)
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
	program, _ := p.ParseProgram()

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
	program, _ := p.ParseProgram()

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
		program, _ := p.ParseProgram()

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

func TestBoolInfix(t *testing.T) {
	tests := []struct {
		input    string
		left     bool
		operator string
		right    bool
	}{
		{"true == true", true, "==", true},
		{"true != true", true, "!=", true},
		// {"true == false", true, "==", false},
		// {"true != false", true, "!=", false},
		// {"false == true", false, "==", true},
		// {"false != true", false, "!=", true},
		// {"false == false", false, "==", false},
		// {"false != false", false, "!=", false},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		err := l.Lex()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		p := NewParser(l)
		program, _ := p.ParseProgram()

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

		left, ok := infix.Left.(*ast.BooleanLiteral)
		if !ok {
			t.Errorf("Expected left to be *ast.BooleanLiteral, got %T", infix.Left)
		}

		if left.Value != tt.left {
			t.Errorf("Expected left to be %t, got %t", tt.left, left.Value)
		}

		right, ok := infix.Right.(*ast.BooleanLiteral)
		if !ok {
			t.Errorf("Expected right to be *ast.BooleanLiteral, got %T", infix.Right)
		}

		if right.Value != tt.right {
			t.Errorf("Expected right to be %t, got %t", tt.right, right.Value)
		}

	}
}

func TestParseFuncLiteral(t *testing.T) {
	source := "@(x, y) { x + y }"
	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program, _ := p.ParseProgram()

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
		t.Errorf("Expected parameter to be 'x', got %s", function.Parameters[0].Str())
	}

	if function.Parameters[1].Str() != "y" {
		t.Errorf("Expected parameter to be 'y', got %s", function.Parameters[1].Str())
	}

	expr := function.Body
	if expr == nil {
		t.Errorf("Expected single-expr body to be not nil")
	}

	if len(expr.Statements) != 1 {
		t.Errorf("Expected 1 statement in body, got %d", len(expr.Statements))
	}

	infix, ok := expr.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.ExpressionStatement, got %T", expr.Statements[0])
	}

	if infix.Expression.Str() != "(x + y)" {
		t.Errorf("Expected expression to be (x + y), got %s", infix.Expression.Str())
	}
}

func TestParseFuncLiteralMultiline(t *testing.T) {
	source := `fn(x, y) {
		z = x + y
		z * 2
	}`

	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program, _ := p.ParseProgram()

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
		t.Errorf("Expected parameter to be 'x', got %s", function.Parameters[0].Str())
	}

	if function.Parameters[1].Str() != "y" {
		t.Errorf("Expected parameter to be 'y', got %s", function.Parameters[1].Str())
	}

	if len(function.Body.Statements) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(function.Body.Statements))
	}

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.LetStatement, got %T", function.Body.Statements[0])
	}

	if body.Expression.Str() != "(z = (x + y))" {
		t.Errorf("Expected expression to be z = (x + y), got %s", body.Expression.Str())
	}

	body, ok = function.Body.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.LetStatement, got %T", function.Body.Statements[1])
	}

	if body.Expression.Str() != "(z * 2)" {
		t.Errorf("Expected expression to be (z * 2), got %s", body.Expression.Str())
	}

}

func TestMultipleStatements(t *testing.T) {
	source := `
	f = fn(x, y) {
		x + y
	}

	g = 100
	`

	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program, _ := p.ParseProgram()

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.LetStatement, got %T", program.Statements[0])
	}

	stmt, ok = program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected statement to be *ast.LetStatement, got %T", program.Statements[1])
	}

	if stmt.Expression.Str() != "(g = 100)" {
		t.Errorf("Expected expression to be g = 100, got %s", stmt.Expression.Str())
	}

}

func TestParseFunctionCall(t *testing.T) {
	source := `f(x, y)`

	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program, _ := p.ParseProgram()

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

	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Errorf("Expected expression to be *ast.CallExpression, got %T", stmt.Expression)
	}

	if call.Function.Str() != "f" {
		t.Errorf("Expected function to be 'f', got %s", call.Function)
	}

	if len(call.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(call.Arguments))
	}

	if call.Arguments[0].Str() != "x" {
		t.Errorf("Expected argument to be 'x', got %s", call.Arguments[0])
	}

	if call.Arguments[1].Str() != "y" {
		t.Errorf("Expected argument to be 'y', got %s", call.Arguments[1])
	}
}

func TestParseTernary(t *testing.T) {
	source := `x ? y : z`

	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program, _ := p.ParseProgram()

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

	ternary, ok := stmt.Expression.(*ast.TernaryExpression)
	if !ok {
		t.Errorf("Expected expression to be *ast.TernaryExpression, got %T", stmt.Expression)
	}

	if ternary.Condition.Str() != "x" {
		t.Errorf("Expected condition to be 'x', got %s", ternary.Condition)
	}

	if ternary.Consequence.Str() != "y" {
		t.Errorf("Expected consequence to be 'y', got %s", ternary.Consequence)
	}

	if ternary.Alternative.Str() != "z" {
		t.Errorf("Expected alternative to be 'z', got %s", ternary.Alternative)
	}
}

func TestParseList(t *testing.T) {
	source := `[1, 2, 3]`

	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p := NewParser(l)
	program, _ := p.ParseProgram()

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

	list, ok := stmt.Expression.(*ast.ListLiteral)
	if !ok {
		t.Errorf("Expected expression to be *ast.ListLiteral, got %T", stmt.Expression)
	}

	if len(list.Elements) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(list.Elements))
	}

	if list.Elements[0].Str() != "1" {
		t.Errorf("Expected element to be '1', got %s", list.Elements[0])
	}

	if list.Elements[1].Str() != "2" {
		t.Errorf("Expected element to be '2', got %s", list.Elements[1])
	}

	if list.Elements[2].Str() != "3" {
		t.Errorf("Expected element to be '3', got %s", list.Elements[2])
	}
}

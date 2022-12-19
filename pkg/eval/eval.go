package eval

import (
	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/parser"
	"github.com/ariaghora/zmol/pkg/val"
)

type Env struct {
	symTable map[string]val.ZValue
}

type ZmolState struct {
	Env *Env
}

func NewZmolState() *ZmolState {
	return &ZmolState{
		Env: &Env{
			symTable: make(map[string]val.ZValue),
		},
	}
}

func (s *ZmolState) Eval(source string) val.ZValue {
	l := lexer.NewLexer(source)
	l.Lex()
	p := parser.NewParser(l)
	program := p.ParseProgram()
	return s.EvalProgram(program)
}

func (s *ZmolState) EvalProgram(node ast.Node) val.ZValue {
	switch node := node.(type) {
	case *ast.Program:
		return s.evalStatements(node)
	case *ast.ExpressionStatement:
		return s.EvalProgram(node.Expression)
	case *ast.InfixExpression:
		left := s.EvalProgram(node.Left)
		right := s.EvalProgram(node.Right)
		return s.evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return s.evalIntegerLiteral(node)
	case *ast.FloatLiteral:
		return s.evalFloatLiteral(node)
	}
	return nil
}

func (s *ZmolState) evalStatements(program *ast.Program) val.ZValue {
	var result val.ZValue

	for _, statement := range program.Statements {
		result = s.EvalProgram(statement)
	}

	return result
}

func (s *ZmolState) evalIntegerLiteral(il *ast.IntegerLiteral) val.ZValue {
	return &val.ZInt{Value: il.Value}
}

func (s *ZmolState) evalFloatLiteral(fl *ast.FloatLiteral) val.ZValue {
	return &val.ZFloat{Value: fl.Value}
}

func (s *ZmolState) evalInfixExpression(operator string, left, right val.ZValue) val.ZValue {
	switch {
	case left.Type() == val.ZINT && right.Type() == val.ZINT:
		return s.evalIntegerInfixExpression(operator, left, right)
	case left.Type() == val.ZFLOAT && right.Type() == val.ZFLOAT:
		return s.evalFloatInfixExpression(operator, left, right)
	case left.Type() == val.ZINT && right.Type() == val.ZFLOAT:
		return s.evalIntFloatInfixExpression(operator, left, right)
	case left.Type() == val.ZFLOAT && right.Type() == val.ZINT:
		return s.evalFloatIntInfixExpression(operator, left, right)
	}
	return nil
}

func (s *ZmolState) evalIntegerInfixExpression(operator string, left, right val.ZValue) val.ZValue {
	leftVal := left.(*val.ZInt).Value
	rightVal := right.(*val.ZInt).Value

	switch operator {
	case "+":
		return &val.ZInt{Value: leftVal + rightVal}
	case "-":
		return &val.ZInt{Value: leftVal - rightVal}
	case "*":
		return &val.ZInt{Value: leftVal * rightVal}
	case "/":
		return &val.ZInt{Value: leftVal / rightVal}
	}
	return nil
}

func (s *ZmolState) evalFloatInfixExpression(operator string, left, right val.ZValue) val.ZValue {
	leftVal := left.(*val.ZFloat).Value
	rightVal := right.(*val.ZFloat).Value

	switch operator {
	case "+":
		return &val.ZFloat{Value: leftVal + rightVal}
	case "-":
		return &val.ZFloat{Value: leftVal - rightVal}
	case "*":
		return &val.ZFloat{Value: leftVal * rightVal}
	case "/":
		return &val.ZFloat{Value: leftVal / rightVal}
	}
	return nil
}

func (s *ZmolState) evalIntFloatInfixExpression(operator string, left, right val.ZValue) val.ZValue {
	leftVal := float64(left.(*val.ZInt).Value)
	rightVal := right.(*val.ZFloat).Value

	switch operator {
	case "+":
		return &val.ZFloat{Value: leftVal + rightVal}
	case "-":
		return &val.ZFloat{Value: leftVal - rightVal}
	case "*":
		return &val.ZFloat{Value: leftVal * rightVal}
	case "/":
		return &val.ZFloat{Value: leftVal / rightVal}
	}
	return nil
}

func (s *ZmolState) evalFloatIntInfixExpression(operator string, left, right val.ZValue) val.ZValue {
	leftVal := left.(*val.ZFloat).Value
	rightVal := float64(right.(*val.ZInt).Value)

	switch operator {
	case "+":
		return &val.ZFloat{Value: leftVal + rightVal}
	case "-":
		return &val.ZFloat{Value: leftVal - rightVal}
	case "*":
		return &val.ZFloat{Value: leftVal * rightVal}
	case "/":
		return &val.ZFloat{Value: leftVal / rightVal}
	}
	return nil
}

package eval

import (
	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/val"
)

func Eval(node ast.Node) val.ZValue {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.FloatLiteral:
		return evalFloatLiteral(node)
	}
	return nil
}

func evalStatements(program *ast.Program) val.ZValue {
	var result val.ZValue

	for _, statement := range program.Statements {
		result = Eval(statement)
	}

	return result
}

func evalIntegerLiteral(il *ast.IntegerLiteral) val.ZValue {
	return &val.ZInt{Value: il.Value}
}

func evalFloatLiteral(fl *ast.FloatLiteral) val.ZValue {
	return &val.ZFloat{Value: fl.Value}
}

func evalInfixExpression(operator string, left, right val.ZValue) val.ZValue {
	switch {
	case left.Type() == val.ZINT && right.Type() == val.ZINT:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == val.ZFLOAT && right.Type() == val.ZFLOAT:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == val.ZINT && right.Type() == val.ZFLOAT:
		return evalIntFloatInfixExpression(operator, left, right)
	case left.Type() == val.ZFLOAT && right.Type() == val.ZINT:
		return evalFloatIntInfixExpression(operator, left, right)
	}
	return nil
}

func evalIntegerInfixExpression(operator string, left, right val.ZValue) val.ZValue {
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

func evalFloatInfixExpression(operator string, left, right val.ZValue) val.ZValue {
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

func evalIntFloatInfixExpression(operator string, left, right val.ZValue) val.ZValue {
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

func evalFloatIntInfixExpression(operator string, left, right val.ZValue) val.ZValue {
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

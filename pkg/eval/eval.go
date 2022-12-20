package eval

import (
	"fmt"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/parser"
	"github.com/ariaghora/zmol/pkg/val"
)

type ZmolState struct {
	Env *val.Env
}

func NewZmolState(ParentEnv *val.Env) *ZmolState {
	symTable := make(map[string]val.ZValue)
	return &ZmolState{
		Env: &val.Env{
			SymTable:  symTable,
			ParentEnv: ParentEnv,
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
	case *ast.Identifier:
		return s.evalIdentifier(node)
	case *ast.VarrAssignmentStatement:
		fmt.Println("WARNING: let is deprecated, omit the 'let' keyword")
		val := s.EvalProgram(node.Value)
		if isErr(val) {
			return val
		}
		return s.Env.Set(node.Name.Value, val)
	case *ast.BlockStatement:
		return s.evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return s.EvalProgram(node.Expression)
	case *ast.InfixExpression:
		switch node.Operator {
		case "=":
			return s.evalVariableAssignment(node)
		case "==", "!=", "<", ">", "<=", ">=":
			return s.evalComparisonExpression(node)
		default:
			left := s.EvalProgram(node.Left)
			right := s.EvalProgram(node.Right)
			return s.evalInfixExpression(node.Operator, left, right)
		}
	case *ast.IntegerLiteral:
		return s.evalIntegerLiteral(node)
	case *ast.FloatLiteral:
		return s.evalFloatLiteral(node)
	case *ast.BooleanLiteral:
		return s.evalBooleanLiteral(node)
	case *ast.ListLiteral:
		return s.evalListLiteral(node)
	case *ast.FuncLiteral:
		params := node.Parameters
		body := node.Body
		return &val.ZFunction{
			Params: params,
			Body:   body,
			Env:    s.Env,
		}
	case *ast.TernaryExpression:
		return s.evalTernaryExpression(node)
	case *ast.CallExpression:
		return s.evalCallExpression(node)
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

func (s *ZmolState) evalBooleanLiteral(bl *ast.BooleanLiteral) val.ZValue {
	return &val.ZBool{Value: bl.Value}
}

func (s *ZmolState) evalListLiteral(ll *ast.ListLiteral) val.ZValue {
	var elements []val.ZValue
	for _, element := range ll.Elements {
		elements = append(elements, s.EvalProgram(element))
	}
	return &val.ZList{Elements: elements}
}

func (s *ZmolState) evalComparisonExpression(node *ast.InfixExpression) val.ZValue {
	left := s.EvalProgram(node.Left)
	right := s.EvalProgram(node.Right)

	switch node.Operator {
	case "==":
		return &val.ZBool{Value: left.Equals(right)}
	case "!=":
		return &val.ZBool{Value: !left.Equals(right)}
	// case "<":
	// 	return &val.ZBool{Value: left.LessThan(right)}
	// case ">":
	// 	return &val.ZBool{Value: left.GreaterThan(right)}
	case "<=":
		return &val.ZBool{Value: left.LessThanEquals(right)}
		// case ">=":
		// 	return &val.ZBool{Value: left.GreaterThanEquals(right)}
	}
	return nil
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

func (s *ZmolState) evalIdentifier(node *ast.Identifier) val.ZValue {
	if val, ok := s.Env.Get(node.Value); ok {
		return val
	}
	return nil
}

func (s *ZmolState) evalVariableAssignment(node *ast.InfixExpression) val.ZValue {
	val := s.EvalProgram(node.Right)
	if isErr(val) {
		return val
	}
	return s.Env.Set(node.Left.(*ast.Identifier).Value, val)
}

func (s *ZmolState) evalBlockStatement(block *ast.BlockStatement) val.ZValue {
	var result val.ZValue
	for _, statement := range block.Statements {
		evaluated := s.EvalProgram(statement)
		result = evaluated
	}
	return result
}

func (s *ZmolState) evalExpressions(exps []ast.Expression) []val.ZValue {
	var result []val.ZValue

	for _, e := range exps {
		evaluated := s.EvalProgram(e)
		if isErr(evaluated) {
			return []val.ZValue{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (s *ZmolState) evalCallExpression(node *ast.CallExpression) val.ZValue {
	function := s.EvalProgram(node.Function)
	if isErr(function) {
		return function
	}
	args := s.evalExpressions(node.Arguments)
	if function.Type() == val.ZNATIVE {
		return function.(*val.ZNativeFunc).Fn(args...)
	}

	params := function.(*val.ZFunction).Params
	zState := NewZmolState(s.Env)
	for i, arg := range args {
		if isErr(arg) {
			return arg
		}
		zState.Env.Set(params[i].Value, arg)
	}
	evaluated := zState.EvalProgram(function.(*val.ZFunction).Body)
	return evaluated
}

func (s *ZmolState) evalTernaryExpression(node *ast.TernaryExpression) val.ZValue {
	condition := s.EvalProgram(node.Condition)
	if isErr(condition) {
		return condition
	}
	if condition.(*val.ZBool).Value {
		return s.EvalProgram(node.Consequence)
	}
	return s.EvalProgram(node.Alternative)
}

func isErr(obj val.ZValue) bool {
	if obj != nil {
		return obj.Type() == val.ZERROR
	}
	return false
}

package eval

import (
	"errors"
	"fmt"
	"os"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/lexer"
	"github.com/ariaghora/zmol/pkg/parser"
	"github.com/ariaghora/zmol/pkg/val"

	"github.com/fatih/color"
)

type ZmolState struct {
	Env *val.Env
}

func NewZmolState(ParentEnv *val.Env) *ZmolState {
	symTable := make(map[string]val.ZValue)
	// default __moddir__ is the current working directory unless
	// the script is imported as a module
	symTable["__moddir__"] = val.STRING(".")
	return &ZmolState{
		Env: &val.Env{
			SymTable:  symTable,
			ParentEnv: ParentEnv,
		},
	}
}

func (s *ZmolState) Eval(source string) (val.ZValue, error) {
	l := lexer.NewLexer(source)
	err := l.Lex()
	if err != nil {
		return nil, err
	}

	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}
	if len(p.Errors()) != 0 {
		s.printParserErrors(os.Stderr, p.Errors())
		return val.ERROR("Parser errors"), errors.New("parser errors")
	}
	return s.EvalProgram(program), nil
}

func (s *ZmolState) printParserErrors(out *os.File, errors []string) {
	for _, msg := range errors {
		color.Red(msg)
	}
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
	case *ast.IfStatement:
		return s.evalIfStatement(node)
	case *ast.ExpressionStatement:
		return s.EvalProgram(node.Expression)
	case *ast.PipelineExpression:
		return s.evalPipelineExpression(node)
	case *ast.InfixExpression:
		switch node.Operator {
		case "=":
			return s.evalVariableAssignment(node)
		case "==", "!=", "<", ">", "<=", ">=":
			return s.evalBooleanExpression(node)
		case "&&", "||":
			return s.evalLogicalExpression(node)
		default:
			left := s.EvalProgram(node.Left)
			right := s.EvalProgram(node.Right)
			return s.evalInfixExpression(node.Operator, left, right)
		}
	case *ast.IntegerLiteral:
		return s.evalIntegerLiteral(node)
	case *ast.FloatLiteral:
		return s.evalFloatLiteral(node)
	case *ast.StringLiteral:
		return &val.ZString{Value: node.Value}
	case *ast.BooleanLiteral:
		return s.evalBooleanLiteral(node)
	case *ast.ListLiteral:
		return s.evalListLiteral(node)
	case *ast.IndexExpression:
		return s.evalIndexExpression(node)
	case *ast.MemberAccessExpression:
		return s.evalMemberAccessExpression(node)
	case *ast.FuncLiteral:
		params := node.Parameters
		body := node.Body
		return val.FUNCTION(params, body, s.Env)
	case *ast.TernaryExpression:
		return s.evalTernaryExpression(node)
	case *ast.IterStatement:
		return s.evalIterStatement(node)
	case *ast.CallExpression:
		return s.evalCallExpression(node)
	}
	RuntimeErrorf("Unknown node type: %T", node)
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

func (s *ZmolState) evalIndexExpression(ie *ast.IndexExpression) val.ZValue {
	left := s.EvalProgram(ie.Left)
	index := s.EvalProgram(ie.Index)

	switch {
	case left.Type() == val.ZLIST && index.Type() == val.ZINT:
		return s.evalListIndexExpression(left, index)
	case left.Type() == val.ZSTRING && index.Type() == val.ZINT:
		return s.evalStringIndexExpression(left, index)
	default:
		return val.ERROR("cannot perform indexing on " + string(left.Type()) + " type")
	}
}

func (s *ZmolState) evalMemberAccessExpression(mae *ast.MemberAccessExpression) val.ZValue {
	left := s.EvalProgram(mae.Left)

	// check if left is a dot accessable type
	leftAccessable, ok := left.(val.ZDotAccessable)
	if !ok {
		return val.ERROR("cannot perform member access on " + string(left.Type()) + " type")
	}

	// check if right is an identifier
	right := mae.Member

	// check if right node is an identifier
	if right.Token.Type != lexer.TokIdent {
		return val.ERROR("cannot perform member access on " + string(right.Token.Type) + " type")
	}

	accessedValue := leftAccessable.DotAccess(right.Value)
	// if it is a function, wrap it as a module function
	if accessedValue.Type() == val.ZFUNCTION {
		return &val.ZModuleFunc{
			Func: accessedValue.(*val.ZFunction),
			Env:  leftAccessable.Env(),
		}
	} else {
		return accessedValue
	}
	// return leftAccessable.DotAccess(right.Value)
}

func (s *ZmolState) evalListIndexExpression(list val.ZValue, index val.ZValue) val.ZValue {
	listVal := list.(*val.ZList)
	indexVal := index.(*val.ZInt)
	max := int64(len(listVal.Elements) - 1)

	if indexVal.Value < 0 || indexVal.Value > max {
		RuntimeErrorf("index out of range: %d", indexVal.Value)
	}

	return listVal.Elements[indexVal.Value]
}

func (s *ZmolState) evalStringIndexExpression(str val.ZValue, index val.ZValue) val.ZValue {
	strVal := str.(*val.ZString)
	indexVal := index.(*val.ZInt)
	max := int64(len(strVal.Value) - 1)

	if indexVal.Value < 0 || indexVal.Value > max {
		RuntimeErrorf("index out of range: %d", indexVal.Value)
	}

	return &val.ZString{Value: string(strVal.Value[indexVal.Value])}
}

func (s *ZmolState) evalMemberAssignment(mae *ast.MemberAccessExpression, value val.ZValue) val.ZValue {
	left := s.EvalProgram(mae.Left)

	// check if left is a dot accessable type
	leftAccessable, ok := left.(val.ZDotAccessable)
	if !ok {
		RuntimeErrorf("cannot perform member access on " + string(left.Type()) + " type")
	}

	// check if right is an identifier
	right := mae.Member

	// check if right node is an identifier
	if right.Token.Type != lexer.TokIdent {
		RuntimeErrorf("identifier expected, got " + string(right.Token.Type) + " type")
	}

	leftAccessable.DotAssign(right.Value, value)
	return val.NULL()
}

func (s *ZmolState) evalIndexAssignment(ie *ast.IndexExpression, value val.ZValue) val.ZValue {
	left := s.EvalProgram(ie.Left)
	index := s.EvalProgram(ie.Index)

	if left.Type() == val.ZLIST && index.Type() == val.ZINT {
		return s.evalListIndexAssignment(left, index, value)
	}
	RuntimeErrorf("index assignment not supported: %s", string(left.Type()))
	return val.NULL()
}

func (s *ZmolState) evalListIndexAssignment(list val.ZValue, index val.ZValue, value val.ZValue) val.ZValue {
	listVal := list.(*val.ZList)
	indexVal := index.(*val.ZInt)
	max := int64(len(listVal.Elements) - 1)

	if indexVal.Value < 0 || indexVal.Value > max {
		RuntimeErrorf("index out of range: %d", indexVal.Value)
	}

	listVal.Elements[indexVal.Value] = value
	return value
}

func (s *ZmolState) evalBooleanExpression(node *ast.InfixExpression) val.ZValue {
	left := s.EvalProgram(node.Left)
	right := s.EvalProgram(node.Right)

	leftComparable, ok := left.(val.ZComparable)
	if !ok {
		return val.ERROR(fmt.Sprintf("cannot compare %s to %s with `%s` operator", left.Type(), right.Type(), node.Operator))
	}

	rightComparable, ok := right.(val.ZComparable)
	if !ok {
		return val.ERROR(fmt.Sprintf("cannot compare %s to %s with `%s` operator", left.Type(), right.Type(), node.Operator))
	}

	switch node.Operator {
	case "==":
		return leftComparable.Equals(rightComparable)
	case "!=":
		return leftComparable.NotEquals(rightComparable)
	case "<":
		return leftComparable.LessThan(rightComparable)
	case ">":
		return leftComparable.GreaterThan(rightComparable)
	case "<=":
		return leftComparable.LessThanEquals(rightComparable)
	case ">=":
		return leftComparable.GreaterThanEquals(rightComparable)
	}
	return val.ERROR(fmt.Sprintf("unknown operator: %s %s %s", node.Left.Str(), node.Operator, node.Right.Str()))
}

func (s *ZmolState) evalLogicalExpression(node *ast.InfixExpression) val.ZValue {
	left := s.EvalProgram(node.Left)
	if isErr(left) {
		return left
	}
	switch node.Operator {
	case "&&":
		if left.(*val.ZBool).Value {
			return s.EvalProgram(node.Right)
		}
		return left
	case "||":
		if !left.(*val.ZBool).Value {
			return s.EvalProgram(node.Right)
		}
		return left
	}
	return val.ERROR(fmt.Sprintf("unknown operator: %s %s %s", node.Left.Str(), node.Operator, node.Right.Str()))
}

func (s *ZmolState) evalPipelineExpression(node *ast.PipelineExpression) val.ZValue {
	list := s.EvalProgram(node.List)
	if isErr(list) {
		return list
	}

	function := s.EvalProgram(node.FuncLiteral)
	if function.Type() != val.ZFUNCTION && function.Type() != val.ZNATIVE && function.Type() != val.ZMODULEFUNC {
		return val.ERROR("Right side of pipeline must be a function")
	}

	extraArgs := []val.ZValue{}
	for _, arg := range node.ExtraArgs {
		extraArgs = append(extraArgs, s.EvalProgram(arg))
	}

	switch node.Token.Type {
	case lexer.TokPipe:
		finalArgs := []val.ZValue{list}
		finalArgs = append(finalArgs, extraArgs...)
		if function.Type() == val.ZNATIVE {
			return function.(*val.ZNativeFunc).Fn(finalArgs...)
		}
		callable, ok := function.(val.ZCallable)
		if !ok {
			return val.ERROR("Right side of pipeline must be a callable, got " + string(function.Type()))
		}

		return s.applyPipe(callable, finalArgs)
	case lexer.TokMap:
		return s.applyMap(function, list, extraArgs)
	case lexer.TokFilter:
		return s.applyFilter(function, list, extraArgs)
	}

	msg := fmt.Sprintf("Unknown pipeline operator: %s", node.Token.Text)
	RuntimeErrorf(msg)
	return val.ERROR(msg)
}

func (s *ZmolState) applyPipe(fn val.ZCallable, args []val.ZValue) val.ZValue {
	if len(args) != len(fn.Params()) {
		fmt.Println("args", args)
		RuntimeErrorf("Wrong number of arguments for `%s`: expected=%d, got=%d", fn.Name(), len(fn.Params()), len(args))
	}
	evaluated := EvalCallable(fn, args, s.Env)

	// zState := NewZmolState(s.Env)

	// for i, param := range fn.Params() {
	// 	zState.Env.Set(param.Value, args[i])
	// }

	// evaluated := zState.EvalProgram(fn.Body())
	return evaluated
}

func (s *ZmolState) applyMap(fn val.ZValue, list val.ZValue, extraArgs []val.ZValue) val.ZValue {
	// error if the list is not a list or string
	if list.Type() != val.ZLIST && list.Type() != val.ZSTRING {
		RuntimeErrorf("Left side of pipeline must be iterable")
	}

	//// Case 1: Native function
	if fn.Type() == val.ZNATIVE {
		newList := &val.ZList{Elements: []val.ZValue{}}
		switch list.Type() {
		case val.ZSTRING:
			for _, elem := range list.(*val.ZString).Value {
				finalArgs := []val.ZValue{&val.ZString{Value: string(elem)}}
				finalArgs = append(finalArgs, extraArgs...)
				newList.Elements = append(newList.Elements, fn.(*val.ZNativeFunc).Fn(finalArgs...))
			}

		case val.ZLIST:
			for _, elem := range list.(*val.ZList).Elements {
				finalArgs := []val.ZValue{elem}
				finalArgs = append(finalArgs, extraArgs...)
				newList.Elements = append(newList.Elements, fn.(*val.ZNativeFunc).Fn(finalArgs...))
			}
		}

		return newList
	}

	//// Case 2: User-defined function
	udFunc, ok := fn.(val.ZCallable)
	if !ok {
		RuntimeErrorf("Right side of pipeline must be a callable, but got %s", fn.Type())
	}
	// zState := NewZmolState(s.Env)
	if len(extraArgs)+1 != len(udFunc.Params()) {
		RuntimeErrorf("Wrong number of arguments: expected=%d, got=%d", len(udFunc.Params()), len(extraArgs)+1)
	}

	// Evaluate the function for each element in the list
	newList := &val.ZList{Elements: []val.ZValue{}}
	switch list.Type() {
	case val.ZSTRING:
		for _, elem := range list.(*val.ZString).Value {
			finalArgs := []val.ZValue{&val.ZString{Value: string(elem)}}
			finalArgs = append(finalArgs, extraArgs...)
			evaluated := EvalCallable(udFunc, finalArgs, s.Env)
			newList.Elements = append(newList.Elements, evaluated)
		}

	case val.ZLIST:
		for _, elem := range list.(*val.ZList).Elements {
			finalArgs := []val.ZValue{elem}
			finalArgs = append(finalArgs, extraArgs...)
			evaluated := EvalCallable(udFunc, finalArgs, s.Env)
			newList.Elements = append(newList.Elements, evaluated)
		}
	}

	return newList
}

func (s *ZmolState) applyFilter(fn val.ZValue, list val.ZValue, extraArgs []val.ZValue) val.ZValue {
	// error if the list is not a list or string
	if list.Type() != val.ZLIST && list.Type() != val.ZSTRING {
		RuntimeErrorf("Left side of pipeline must be iterable")
	}

	//// Case 1: Native function
	if fn.Type() == val.ZNATIVE {
		newList := &val.ZList{Elements: []val.ZValue{}}
		switch list.Type() {
		case val.ZSTRING:
			for _, elem := range list.(*val.ZString).Value {
				finalArgs := []val.ZValue{&val.ZString{Value: string(elem)}}
				finalArgs = append(finalArgs, extraArgs...)
				if fn.(*val.ZNativeFunc).Fn(finalArgs...).(*val.ZBool).Value {
					newList.Elements = append(newList.Elements, &val.ZString{Value: string(elem)})
				}
			}

		case val.ZLIST:
			for _, elem := range list.(*val.ZList).Elements {
				finalArgs := []val.ZValue{elem}
				finalArgs = append(finalArgs, extraArgs...)
				if fn.(*val.ZNativeFunc).Fn(finalArgs...).(*val.ZBool).Value {
					newList.Elements = append(newList.Elements, elem)
				}
			}
		}

		return newList
	}

	//// Case 2: User-defined function
	udFunc, ok := fn.(val.ZCallable)
	if !ok {
		RuntimeErrorf("Right side of pipeline must be a callable, but got %s", fn.Type())
	}
	// zState := NewZmolState(s.Env)
	if len(extraArgs)+1 != len(udFunc.Params()) {
		RuntimeErrorf("Wrong number of arguments: expected=%d, got=%d", len(udFunc.Params()), len(extraArgs)+1)
	}

	// Evaluate the function for each element in the list
	newList := &val.ZList{Elements: []val.ZValue{}}
	switch list.Type() {
	case val.ZSTRING:
		for _, elem := range list.(*val.ZString).Value {
			finalArgs := []val.ZValue{&val.ZString{Value: string(elem)}}
			finalArgs = append(finalArgs, extraArgs...)
			evaluated := EvalCallable(udFunc, finalArgs, s.Env)
			if evaluated.(*val.ZBool).Value {
				newList.Elements = append(newList.Elements, &val.ZString{Value: string(elem)})
			}

		}
	case val.ZLIST:
		for _, elem := range list.(*val.ZList).Elements {
			finalArgs := []val.ZValue{elem}
			finalArgs = append(finalArgs, extraArgs...)
			evaluated := EvalCallable(udFunc, finalArgs, s.Env)
			if evaluated.(*val.ZBool).Value {
				newList.Elements = append(newList.Elements, elem)
			}
		}
	}

	return newList
}

func EvalCallable(fn val.ZCallable, args []val.ZValue, parentEnv *val.Env) val.ZValue {
	zState := NewZmolState(parentEnv)
	if len(args) != len(fn.Params()) {
		RuntimeErrorf("Wrong number of arguments: expected=%d, got=%d", len(fn.Params()), len(args))
	}

	for i, param := range fn.Params() {
		zState.Env.Set(param.Value, args[i])
	}

	return zState.EvalProgram(fn.Body())
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

	case left.Type() == val.ZLIST && right.Type() == val.ZLIST:
		return s.evalListConcatExpression(left, right)

	case left.Type() == val.ZSTRING && right.Type() == val.ZSTRING:
		return &val.ZString{Value: left.(*val.ZString).Value + right.(*val.ZString).Value}
	}

	return val.ERROR(fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type()))
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
	case "%":
		return &val.ZInt{Value: leftVal % rightVal}
	}
	return val.ERROR("Operator " + operator + " not supported for integers")
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
	return val.ERROR("Operator " + operator + " not supported for floats")
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
	return val.ERROR("Operator " + operator + " not supported for integers and floats")
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
	return val.ERROR("Operator " + operator + " not supported for floats and integers")
}

func (s *ZmolState) evalIdentifier(node *ast.Identifier) val.ZValue {
	if val, ok := s.Env.Get(node.Value); ok {
		return val
	}
	return RuntimeErrorf("identifier not found: " + node.Value)
}

func RuntimeErrorf(format string, args ...interface{}) val.ZValue {
	color.Set(color.FgRed)
	fmt.Println("\n*** RUNTIME ERROR ***")
	fmt.Println(fmt.Sprintf(format, args...))
	color.Unset()
	os.Exit(1)
	return nil
}

func (s *ZmolState) evalVariableAssignment(node *ast.InfixExpression) val.ZValue {
	value := s.EvalProgram(node.Right)
	if isErr(value) {
		return value
	}

	switch node.Left.(type) {
	case *ast.MemberAccessExpression:
		return s.evalMemberAssignment(node.Left.(*ast.MemberAccessExpression), value)
	case *ast.Identifier:
		if value.Type() == val.ZCLASS {
			value.(*val.ZClass).Name = node.Left.(*ast.Identifier).Value
		} else if value.Type() == val.ZFUNCTION {
			value.(*val.ZFunction).SetName(node.Left.(*ast.Identifier).Value)
		}

		return s.Env.Set(node.Left.(*ast.Identifier).Value, value)
	case *ast.IndexExpression:
		return s.evalIndexAssignment(node.Left.(*ast.IndexExpression), value)
	}

	RuntimeErrorf("invalid assignment")
	return val.NULL()
}

func (s *ZmolState) evalBlockStatement(block *ast.BlockStatement) val.ZValue {
	var result val.ZValue
	for _, statement := range block.Statements {
		evaluated := s.EvalProgram(statement)
		result = evaluated
	}
	return result
}

func (s *ZmolState) evalIfStatement(node *ast.IfStatement) val.ZValue {
	condition := s.EvalProgram(node.Condition)
	if isErr(condition) {
		return condition
	}
	if isTruthy(condition) {
		return s.EvalProgram(node.Consequence)
	} else if node.Alternative != nil {
		return s.EvalProgram(node.Alternative)
	}
	return val.NULL()
}

func isTruthy(value val.ZValue) bool {
	switch value.Type() {
	case val.ZNULL:
		return false
	case val.ZBOOL:
		return value.(*val.ZBool).Value
	default:
		return true
	}
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
	lEval := s.EvalProgram(node.Function)
	if isErr(lEval) {
		return lEval
	}

	args := s.evalExpressions(node.Arguments)

	if lEval.Type() == val.ZNATIVE {
		function := lEval.(*val.ZNativeFunc)
		return function.Fn(args...)
	} else if lEval.Type() == val.ZMODULEFUNC {
		function := lEval
		params := function.(*val.ZModuleFunc).Func.Params()
		if len(args) == 1 && isErr(args[0]) {
			return args[0]
		}
		zState := NewZmolState(function.(*val.ZModuleFunc).Env)
		for i, arg := range args {
			if isErr(arg) {
				return arg
			}
			zState.Env.Set(params[i].Value, arg)
		}
		evaluated := zState.EvalProgram(function.(*val.ZModuleFunc).Func.Body())
		return evaluated
	} else if lEval.Type() == val.ZCLASS {
		class := lEval.(*val.ZClass)
		obj := val.OBJECT(class.Name, s.Env)

		// Copy class attributes and methods (i.e., ZValues in the class' env)
		// to the new object
		for k, v := range class.Env().SymTable {
			obj.Env().Set(k, v)
		}
		obj.Env().ParentEnv = class.Env().ParentEnv

		// try to call the constructor. If it doesn't exist, just return the new object
		constructorFn, ok := obj.Env().Get("init")
		if ok {
			// if the constructor exists, call it
			params := constructorFn.(*val.ZFunction).Params()

			if len(args) != len(params) {
				RuntimeErrorf("wrong number of arguments for constructor. got=%d, want=%d", len(args), len(params))
			}

			zState := NewZmolState(obj.Env())
			for i, arg := range args {
				if isErr(arg) {
					return arg
				}
				zState.Env.Set(params[i].Value, arg)
			}
			zState.EvalProgram(constructorFn.(*val.ZFunction).Body())
		}

		if !ok && len(args) > 0 {
			RuntimeErrorf("wrong number of arguments for constructor. got=%d, want=0", len(args))
		}

		return obj

	} else {
		function := s.EvalProgram(node.Function)
		if isErr(function) {
			return function
		}

		env := s.Env
		params := function.(*val.ZFunction).Params()
		zState := NewZmolState(env)
		for i, arg := range args {
			if isErr(arg) {
				return arg
			}
			zState.Env.Set(params[i].Value, arg)
		}
		evaluated := zState.EvalProgram(function.(*val.ZFunction).Body())
		return evaluated
	}

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

func (s *ZmolState) evalIterStatement(node *ast.IterStatement) val.ZValue {
	list := s.EvalProgram(node.List)

	if isErr(list) {
		return list
	}

	// check if it is a list
	if list.Type() != val.ZLIST {
		fmt.Println("Iter statement requires a list")
		os.Exit(1)
	}

	ident := node.Ident.Value

	for _, item := range list.(*val.ZList).Elements {
		s.Env.Set(ident, item)
		s.EvalProgram(node.Body)
	}

	return val.NULL()
}

func (s *ZmolState) evalListConcatExpression(left, right val.ZValue) val.ZValue {
	if isErr(left) {
		return left
	}
	if isErr(right) {
		return right
	}

	// check if both are lists
	if left.Type() != val.ZLIST || right.Type() != val.ZLIST {
		fmt.Println("Concatenation requires two lists")
		os.Exit(1)
	}

	return &val.ZList{Elements: append(left.(*val.ZList).Elements, right.(*val.ZList).Elements...)}
}

func isErr(obj val.ZValue) bool {
	if obj != nil {
		return obj.Type() == val.ZERROR
	}
	return false
}

package val

import (
	"github.com/ariaghora/zmol/pkg/ast"
)

// Function type
type ZFunction struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Env
}

func (z *ZFunction) Type() ZValueType { return ZFUNCTION }
func (z *ZFunction) Str() string      { return "Function" }
func (z *ZFunction) Equals(other ZValue) ZValue {
	return ERROR("Operator '==' not defined for functions")
}
func (z *ZFunction) NotEquals(other ZValue) ZValue {
	return ERROR("Operator '!=' not defined for functions")
}
func (z *ZFunction) LessThanEquals(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for functions")
}
func (z *ZFunction) GreaterThanEquals(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for functions")
}

// Native built-in function type
type ZNativeFunc struct {
	Fn func(args ...ZValue) ZValue
}

func (z *ZNativeFunc) Type() ZValueType { return ZNATIVE }
func (z *ZNativeFunc) Str() string      { return "NativeFunction" }
func (z *ZNativeFunc) Equals(other ZValue) ZValue {
	return ERROR("Operator '==' not defined for functions")
}
func (z *ZNativeFunc) NotEquals(other ZValue) ZValue {
	return ERROR("Operator '!=' not defined for functions")
}
func (z *ZNativeFunc) LessThanEquals(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for functions")
}
func (z *ZNativeFunc) GreaterThanEquals(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for functions")
}

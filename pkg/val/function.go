package val

import "github.com/ariaghora/zmol/pkg/ast"

// Function type
type ZFunction struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Env
}

func (z *ZFunction) Type() ZValueType { return ZFUNCTION }
func (z *ZFunction) Str() string      { return "Function" }
func (z *ZFunction) Equals(other ZValue) bool {
	if other.Type() != ZFUNCTION {
		return false
	}
	return false
}
func (z *ZFunction) LessThanEquals(other ZValue) bool {
	return false
}

// Native built-in function type
type ZNativeFunc struct {
	Fn func(args ...ZValue) ZValue
}

func (z *ZNativeFunc) Type() ZValueType { return ZNATIVE }
func (z *ZNativeFunc) Str() string      { return "NativeFunction" }
func (z *ZNativeFunc) Equals(other ZValue) bool {
	if other.Type() != ZNATIVE {
		return false
	}
	return false
}
func (z *ZNativeFunc) LessThanEquals(other ZValue) bool {
	return false
}

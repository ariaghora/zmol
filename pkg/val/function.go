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
func (z *ZFunction) Call(args ...ZValue) ZValue {
	return nil
}

// Native built-in function type
type ZNativeFunc struct {
	Fn func(args ...ZValue) ZValue
}

func (z *ZNativeFunc) Type() ZValueType { return ZNATIVE }
func (z *ZNativeFunc) Str() string      { return "NativeFunction" }
func (z *ZNativeFunc) Call(args ...ZValue) ZValue {
	return z.Fn(args...)
}

type ZModuleFunc struct {
	Func *ZFunction
	Env  *Env
}

func (z *ZModuleFunc) Type() ZValueType { return ZMODULEFUNC }
func (z *ZModuleFunc) Str() string      { return "ModuleFunction" }
func (z *ZModuleFunc) Call(args ...ZValue) ZValue {
	// return z.Func.Call(args...)
	return nil
}

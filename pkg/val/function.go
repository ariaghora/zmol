package val

import (
	"github.com/ariaghora/zmol/pkg/ast"
)

// Function type
type ZFunction struct {
	name   string
	params []*ast.Identifier
	body   *ast.BlockStatement
	Env    *Env
}

func FUNCTION(params []*ast.Identifier, body *ast.BlockStatement, env *Env) *ZFunction {
	return &ZFunction{"<anonymous function>", params, body, env}
}

func (z *ZFunction) Type() ZValueType { return ZFUNCTION }
func (z *ZFunction) Str() string      { return "Function" }
func (z *ZFunction) Call(args ...ZValue) ZValue {
	return nil
}
func (z *ZFunction) Body() *ast.BlockStatement {
	return z.body
}
func (z *ZFunction) Params() []*ast.Identifier {
	return z.params
}
func (z *ZFunction) Name() string {
	return z.name
}
func (z *ZFunction) SetName(name string) {
	z.name = name
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
func (z *ZModuleFunc) Body() *ast.BlockStatement {
	return z.Func.Body()
}
func (z *ZModuleFunc) Params() []*ast.Identifier {
	return z.Func.Params()
}
func (z *ZModuleFunc) Name() string {
	return z.Func.Name()
}

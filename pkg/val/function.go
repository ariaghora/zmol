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

// Native built-in function type
type ZNativeFunc struct {
	Fn func(args ...ZValue) ZValue
}

func (z *ZNativeFunc) Type() ZValueType { return ZNATIVE }
func (z *ZNativeFunc) Str() string      { return "NativeFunction" }

type ZModuleFunc struct {
	Func *ZFunction
	Env  *Env
}

func (z *ZModuleFunc) Type() ZValueType { return ZMODULEFUNC }
func (z *ZModuleFunc) Str() string      { return "ModuleFunction" }

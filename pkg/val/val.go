package val

import "github.com/ariaghora/zmol/pkg/ast"

type ZValueType string

const (
	ZINT        ZValueType = "Int"
	ZFLOAT      ZValueType = "Float"
	ZBOOL       ZValueType = "Bool"
	ZLIST       ZValueType = "List"
	ZCLASS      ZValueType = "Class"
	ZOBJECT     ZValueType = "Object"
	ZTENSOR     ZValueType = "Tensor"
	ZERROR      ZValueType = "Error"
	ZSTRING     ZValueType = "String"
	ZFUNCTION   ZValueType = "Function"
	ZNATIVE     ZValueType = "BuiltinFunction"
	ZNULL       ZValueType = "Null"
	ZMODULE     ZValueType = "Module"
	ZMODULEFUNC ZValueType = "ModuleFunction"
)

type Env struct {
	SymTable  map[string]ZValue
	ParentEnv *Env
}

func (e *Env) Get(name string) (ZValue, bool) {
	obj, ok := e.SymTable[name]
	if !ok && e.ParentEnv != nil {
		return e.ParentEnv.Get(name)
	}
	return obj, ok
}

func (e *Env) Set(name string, val ZValue) ZValue {
	e.SymTable[name] = val
	return val
}

// The ZValue interface is implemented by all types that can be used as values
type ZValue interface {
	Type() ZValueType
	Str() string
}

type ZComparable interface {
	ZValue

	Equal(other ZValue) ZValue
	NotEqual(other ZValue) ZValue
	LessThan(other ZValue) ZValue
	GreaterThan(other ZValue) ZValue
	LessThanEqual(other ZValue) ZValue
	GreaterThanEqual(other ZValue) ZValue
}

type ZLogical interface {
	ZValue

	And(other ZValue) ZValue
	Or(other ZValue) ZValue
	Not() ZValue
}

type ZArithOperand interface {
	ZValue
	Add(other ZValue) ZValue
	Sub(other ZValue) ZValue
	Mul(other ZValue) ZValue
	Div(other ZValue) ZValue
	Mod(other ZValue) ZValue
	Neg() ZValue
}

type ZDotAccessable interface {
	ZValue
	DotAccess(name string) ZValue
	DotAssign(name string, value ZValue)
	Env() *Env
}

type ZCallable interface {
	ZValue
	Params() []*ast.Identifier
	Body() *ast.BlockStatement
	Name() string
}

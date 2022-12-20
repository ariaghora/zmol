package val

type ZValueType string

const (
	ZINT      ZValueType = "Int"
	ZFLOAT    ZValueType = "Float"
	ZBOOL     ZValueType = "Bool"
	ZLIST     ZValueType = "List"
	ZERROR    ZValueType = "Error"
	ZSTRING   ZValueType = "String"
	ZFUNCTION ZValueType = "Function"
	ZNATIVE   ZValueType = "BuiltinFunction"
	ZNULL     ZValueType = "Null"
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
	Equals(other ZValue) ZValue
	LessThanEquals(other ZValue) ZValue
	GreaterThanEquals(other ZValue) ZValue
}

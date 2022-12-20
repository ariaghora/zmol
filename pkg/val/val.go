package val

type ZValueType string

const (
	ZINT      ZValueType = "int"
	ZFLOAT    ZValueType = "float"
	ZBOOL     ZValueType = "bool"
	ZERROR    ZValueType = "error"
	ZSTRING   ZValueType = "string"
	ZFUNCTION ZValueType = "function"
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
	Equals(other ZValue) bool
	LessThanEquals(other ZValue) bool
}

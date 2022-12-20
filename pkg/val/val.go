package val

import (
	"fmt"

	"github.com/ariaghora/zmol/pkg/ast"
)

type ZValueType string

const (
	ZINT   ZValueType = "int"
	ZFLOAT ZValueType = "float"
	ZERROR ZValueType = "error"
)

type Env struct {
	SymTable map[string]ZValue
}

func (e *Env) Get(name string) (ZValue, bool) {
	obj, ok := e.SymTable[name]
	return obj, ok
}

func (e *Env) Set(name string, val ZValue) ZValue {
	e.SymTable[name] = val
	return val
}

type ZValue interface {
	Type() ZValueType
	Str() string
}

type ZInt struct {
	Value int64
}

func (z *ZInt) Type() ZValueType { return "int" }
func (z *ZInt) Str() string {
	return fmt.Sprintf("%d", z.Value)
}

type ZFloat struct {
	Value float64
}

func (z *ZFloat) Type() ZValueType { return "float" }
func (z *ZFloat) Str() string {
	return fmt.Sprintf("%f", z.Value)
}

type ZString struct {
	Value string
}

func (z *ZString) Type() ZValueType { return "string" }
func (z *ZString) Str() string      { return z.Value }

type ZError struct {
	Message string
}

func (z *ZError) Type() ZValueType { return "error" }
func (z *ZError) Str() string      { return "ERROR: " + z.Message }

type ZFunction struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Env
}

func (z *ZFunction) Type() ZValueType { return "function" }
func (z *ZFunction) Str() string      { return "function" }

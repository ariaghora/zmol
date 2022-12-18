package val

import "fmt"

type ZValueType string

const (
	ZINT   ZValueType = "int"
	ZFLOAT ZValueType = "float"
)

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

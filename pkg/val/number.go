package val

import (
	"fmt"
	"math"
)

///////////////////////////////////////////
//              Integer type             //
///////////////////////////////////////////

type ZInt struct {
	Value int64
}

func INT(value int64) *ZInt {
	return &ZInt{Value: value}
}

func (z *ZInt) Type() ZValueType { return ZINT }
func (z *ZInt) Str() string {
	return fmt.Sprintf("%d", z.Value)
}

func (z *ZInt) Add(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return INT(z.Value + other.(*ZInt).Value)
	case ZFLOAT:
		return FLOAT(float64(z.Value) + other.(*ZFloat).Value)
	default:
		return OP_ERROR("+", z, other)
	}
}

func (z *ZInt) Sub(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return INT(z.Value - other.(*ZInt).Value)
	case ZFLOAT:
		return FLOAT(float64(z.Value) - other.(*ZFloat).Value)
	default:
		return OP_ERROR("-", z, other)
	}
}

func (z *ZInt) Mul(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return INT(z.Value * other.(*ZInt).Value)
	case ZFLOAT:
		return FLOAT(float64(z.Value) * other.(*ZFloat).Value)
	default:
		return OP_ERROR("*", z, other)
	}
}

func (z *ZInt) Div(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return INT(z.Value / other.(*ZInt).Value)
	case ZFLOAT:
		return FLOAT(float64(z.Value) / other.(*ZFloat).Value)
	default:
		return OP_ERROR("/", z, other)
	}
}

func (z *ZInt) Mod(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return INT(z.Value % other.(*ZInt).Value)
	default:
		return OP_ERROR("%", z, other)
	}
}

func (z *ZInt) Neg() ZValue {
	return INT(-z.Value)
}

func (z *ZInt) Equal(other ZValue) ZValue {
	if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) == other.(*ZFloat).Value}
	} else if other.Type() == ZINT {
		return &ZBool{Value: z.Value == other.(*ZInt).Value}
	}
	return ERROR(fmt.Sprintf("Operator '==' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZInt) NotEqual(other ZValue) ZValue {
	if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) != other.(*ZFloat).Value}
	} else if other.Type() == ZINT {
		return &ZBool{Value: z.Value != other.(*ZInt).Value}
	}
	return ERROR(fmt.Sprintf("Operator '!=' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZInt) LessThan(other ZValue) ZValue {
	if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) < other.(*ZFloat).Value}
	} else if other.Type() == ZINT {
		return &ZBool{Value: z.Value < other.(*ZInt).Value}
	}
	return ERROR(fmt.Sprintf("Operator '<' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZInt) GreaterThan(other ZValue) ZValue {
	if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) > other.(*ZFloat).Value}
	} else if other.Type() == ZINT {
		return &ZBool{Value: z.Value > other.(*ZInt).Value}
	}
	return ERROR(fmt.Sprintf("Operator '>' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZInt) LessThanEqual(other ZValue) ZValue {
	if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) <= other.(*ZFloat).Value}
	} else if other.Type() == ZINT {
		return &ZBool{Value: z.Value <= other.(*ZInt).Value}
	}
	return ERROR(fmt.Sprintf("Operator '<=' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZInt) GreaterThanEqual(other ZValue) ZValue {
	if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) >= other.(*ZFloat).Value}
	} else if other.Type() == ZINT {
		return &ZBool{Value: z.Value >= other.(*ZInt).Value}
	}
	return ERROR(fmt.Sprintf("Operator '>=' not defined for %s and %s", z.Type(), other.Type()))
}

///////////////////////////////////////////
//               Float type              //
///////////////////////////////////////////

type ZFloat struct {
	Value float64
}

func FLOAT(value float64) *ZFloat {
	return &ZFloat{Value: value}
}

func (z *ZFloat) Type() ZValueType { return ZFLOAT }
func (z *ZFloat) Str() string {
	return fmt.Sprintf("%f", z.Value)
}

func (z *ZFloat) Add(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return FLOAT(z.Value + float64(other.(*ZInt).Value))
	case ZFLOAT:
		return FLOAT(z.Value + other.(*ZFloat).Value)
	default:
		return OP_ERROR("+", z, other)
	}
}

func (z *ZFloat) Sub(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return FLOAT(z.Value - float64(other.(*ZInt).Value))
	case ZFLOAT:
		return FLOAT(z.Value - other.(*ZFloat).Value)
	default:
		return OP_ERROR("-", z, other)
	}
}

func (z *ZFloat) Mul(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return FLOAT(z.Value * float64(other.(*ZInt).Value))
	case ZFLOAT:
		return FLOAT(z.Value * other.(*ZFloat).Value)
	default:
		return OP_ERROR("*", z, other)
	}
}

func (z *ZFloat) Div(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return FLOAT(z.Value / float64(other.(*ZInt).Value))
	case ZFLOAT:
		return FLOAT(z.Value / other.(*ZFloat).Value)
	default:
		return OP_ERROR("/", z, other)
	}
}

func (z *ZFloat) Mod(other ZValue) ZValue {
	switch other.Type() {
	case ZINT:
		return FLOAT(math.Mod(z.Value, float64(other.(*ZInt).Value)))
	case ZFLOAT:
		return FLOAT(math.Mod(z.Value, other.(*ZFloat).Value))
	default:
		return OP_ERROR("%", z, other)
	}
}

func (z *ZFloat) Neg() ZValue {
	return FLOAT(-z.Value)
}

func (z *ZFloat) Equals(other ZValue) ZValue {
	if other.Type() == ZINT {
		return BOOL(z.Value == float64(other.(*ZInt).Value))
	} else if other.Type() == ZFLOAT {
		return BOOL(z.Value == other.(*ZFloat).Value)
	}
	return ERROR(fmt.Sprintf("Operator '==' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZFloat) NotEquals(other ZValue) ZValue {
	if other.Type() == ZINT {
		return BOOL(z.Value != float64(other.(*ZInt).Value))
	} else if other.Type() == ZFLOAT {
		return BOOL(z.Value != other.(*ZFloat).Value)
	}
	return ERROR(fmt.Sprintf("Operator '!=' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZFloat) LessThan(other ZValue) ZValue {
	if other.Type() == ZINT {
		return BOOL(z.Value < float64(other.(*ZInt).Value))
	} else if other.Type() == ZFLOAT {
		return BOOL(z.Value < other.(*ZFloat).Value)
	}
	return ERROR(fmt.Sprintf("Operator '<' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZFloat) GreaterThan(other ZValue) ZValue {
	if other.Type() == ZINT {
		return BOOL(z.Value > float64(other.(*ZInt).Value))
	} else if other.Type() == ZFLOAT {
		return BOOL(z.Value > other.(*ZFloat).Value)
	}
	return ERROR(fmt.Sprintf("Operator '>' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZFloat) LessThanEqual(other ZValue) ZValue {
	if other.Type() == ZINT {
		return BOOL(z.Value <= float64(other.(*ZInt).Value))
	} else if other.Type() == ZFLOAT {
		return BOOL(z.Value <= other.(*ZFloat).Value)
	}
	return ERROR(fmt.Sprintf("Operator '<=' not defined for %s and %s", z.Type(), other.Type()))
}

func (z *ZFloat) GreaterThanEqual(other ZValue) ZValue {
	if other.Type() == ZINT {
		return BOOL(z.Value >= float64(other.(*ZInt).Value))
	} else if other.Type() == ZFLOAT {
		return BOOL(z.Value >= other.(*ZFloat).Value)
	}
	return ERROR(fmt.Sprintf("Operator '>=' not defined for %s and %s", z.Type(), other.Type()))
}

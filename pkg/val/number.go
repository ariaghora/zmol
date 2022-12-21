package val

import "fmt"

// Integer type
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
func (z *ZInt) Equals(other ZValue) ZValue {
	if other.Type() != ZINT {
		return &ZBool{Value: false}
	} else if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) == other.(*ZFloat).Value}
	}
	return &ZBool{Value: z.Value == other.(*ZInt).Value}
}

func (z *ZInt) LessThanEquals(other ZValue) ZValue {
	if other.Type() != ZINT {
		return &ZBool{Value: false}
	} else if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) <= other.(*ZFloat).Value}
	}
	return &ZBool{Value: z.Value <= other.(*ZInt).Value}
}

func (z *ZInt) GreaterThanEquals(other ZValue) ZValue {
	if other.Type() != ZINT {
		return &ZBool{Value: false}
	} else if other.Type() == ZFLOAT {
		return &ZBool{Value: float64(z.Value) >= other.(*ZFloat).Value}
	}
	return &ZBool{Value: z.Value >= other.(*ZInt).Value}
}

// Float type
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
func (z *ZFloat) Equals(other ZValue) ZValue {
	if other.Type() != ZFLOAT {
		return BOOL(false)
	} else if other.Type() == ZINT {
		return BOOL(z.Value == float64(other.(*ZInt).Value))
	}
	return BOOL(z.Value == other.(*ZFloat).Value)
}

func (z *ZFloat) LessThanEquals(other ZValue) ZValue {
	if other.Type() != ZFLOAT {
		return BOOL(false)
	} else if other.Type() == ZINT {
		return BOOL(z.Value <= float64(other.(*ZInt).Value))
	}
	return BOOL(z.Value <= other.(*ZFloat).Value)
}

func (z *ZFloat) GreaterThanEquals(other ZValue) ZValue {
	if other.Type() != ZFLOAT {
		return BOOL(false)
	} else if other.Type() == ZINT {
		return BOOL(z.Value >= float64(other.(*ZInt).Value))
	}
	return BOOL(z.Value >= other.(*ZFloat).Value)
}

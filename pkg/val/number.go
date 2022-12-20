package val

import "fmt"

// Integer type
type ZInt struct {
	Value int64
}

func (z *ZInt) Type() ZValueType { return ZINT }
func (z *ZInt) Str() string {
	return fmt.Sprintf("%d", z.Value)
}
func (z *ZInt) Equals(other ZValue) bool {
	if other.Type() != ZINT {
		return false
	} else if other.Type() == ZFLOAT {
		return float64(z.Value) == other.(*ZFloat).Value
	}
	return z.Value == other.(*ZInt).Value
}

func (z *ZInt) LessThanEquals(other ZValue) bool {
	if other.Type() != ZINT {
		return false
	} else if other.Type() == ZFLOAT {
		return float64(z.Value) <= other.(*ZFloat).Value
	}
	return z.Value <= other.(*ZInt).Value
}

// Float type
type ZFloat struct {
	Value float64
}

func (z *ZFloat) Type() ZValueType { return ZFLOAT }
func (z *ZFloat) Str() string {
	return fmt.Sprintf("%f", z.Value)
}
func (z *ZFloat) Equals(other ZValue) bool {
	if other.Type() != ZFLOAT {
		return false
	} else if other.Type() == ZINT {
		return z.Value == float64(other.(*ZInt).Value)
	}
	return z.Value == other.(*ZFloat).Value
}

func (z *ZFloat) LessThanEquals(other ZValue) bool {
	if other.Type() != ZFLOAT {
		return false
	} else if other.Type() == ZINT {
		return z.Value <= float64(other.(*ZInt).Value)
	}
	return z.Value <= other.(*ZFloat).Value
}

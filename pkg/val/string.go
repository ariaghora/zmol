package val

// String type
type ZString struct {
	Value string
}

func STRING(value string) *ZString {
	return &ZString{Value: value}
}

func (z *ZString) Type() ZValueType { return ZSTRING }
func (z *ZString) Str() string      { return z.Value }
func (z *ZString) Equals(other ZValue) ZValue {
	if other.Type() != ZSTRING {
		return BOOL(false)
	}
	return BOOL(z.Value == other.(*ZString).Value)
}
func (z *ZString) NotEquals(other ZValue) ZValue {
	if other.Type() != ZSTRING {
		return BOOL(true)
	}
	return BOOL(z.Value != other.(*ZString).Value)
}
func (z *ZString) LessThanEquals(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for string")
}

func (z *ZString) GreaterThanEquals(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for string")
}

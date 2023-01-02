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
func (z *ZString) Equal(other ZValue) ZValue {
	if other.Type() != ZSTRING {
		return BOOL(false)
	}
	return BOOL(z.Value == other.(*ZString).Value)
}
func (z *ZString) NotEqual(other ZValue) ZValue {
	if other.Type() != ZSTRING {
		return BOOL(true)
	}
	return BOOL(z.Value != other.(*ZString).Value)
}

func (z *ZString) LessThan(other ZValue) ZValue {
	return ERROR("Operator '<' not defined for string")
}

func (z *ZString) GreaterThan(other ZValue) ZValue {
	return ERROR("Operator '>' not defined for string")
}

func (z *ZString) LessThanEqual(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for string")
}

func (z *ZString) GreaterThanEqual(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for string")
}

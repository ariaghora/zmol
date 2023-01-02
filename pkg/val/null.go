package val

type ZNull struct{}

func NULL() *ZNull {
	return &ZNull{}
}

func (z *ZNull) Type() ZValueType { return ZNULL }
func (z *ZNull) Str() string      { return "" }
func (z *ZNull) Equal(other ZValue) ZValue {
	return BOOL(other.Type() == ZNULL)
}
func (z *ZNull) NotEqual(other ZValue) ZValue {
	return BOOL(other.Type() != ZNULL)
}

func (z *ZNull) LessThan(other ZValue) ZValue {
	return ERROR("Operator '<' not defined for null")
}

func (z *ZNull) GreaterThan(other ZValue) ZValue {
	return ERROR("Operator '>' not defined for null")
}

func (z *ZNull) LessThanEqual(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for null")
}

func (z *ZNull) GreaterThanEqual(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for null")
}

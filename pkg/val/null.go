package val

type ZNull struct{}

func (z *ZNull) Type() ZValueType { return ZNULL }
func (z *ZNull) Str() string      { return "" }
func (z *ZNull) Equals(other ZValue) ZValue {
	return BOOL(other.Type() == ZNULL)
}
func (z *ZNull) LessThanEquals(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for null")
}
func (z *ZNull) GreaterThanEquals(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for null")
}

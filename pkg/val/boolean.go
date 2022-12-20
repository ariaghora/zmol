package val

// Bool type
type ZBool struct {
	Value bool
}

func BOOL(value bool) *ZBool {
	return &ZBool{Value: value}
}

func (z *ZBool) Type() ZValueType { return ZBOOL }
func (z *ZBool) Str() string {
	if z.Value {
		return "true"
	}
	return "false"
}
func (z *ZBool) Equals(other ZValue) ZValue {
	if other.Type() != ZBOOL {
		return ERROR("Cannot compare bool with " + string(other.Type()))
	}
	return BOOL(z.Value == other.(*ZBool).Value)
}
func (z *ZBool) LessThanEquals(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for bool")
}
func (z *ZBool) GreaterThanEquals(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for bool")
}

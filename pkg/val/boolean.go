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
func (z *ZBool) Equal(other ZValue) ZValue {
	if other.Type() != ZBOOL {
		return ERROR("Cannot compare bool with " + string(other.Type()))
	}
	return BOOL(z.Value == other.(*ZBool).Value)
}
func (z *ZBool) NotEqual(other ZValue) ZValue {
	if other.Type() != ZBOOL {
		return ERROR("Cannot compare bool with " + string(other.Type()))
	}
	return BOOL(z.Value != other.(*ZBool).Value)
}
func (z *ZBool) LessThan(other ZValue) ZValue {
	return ERROR("Operator '<' not defined for bool")
}
func (z *ZBool) GreaterThan(other ZValue) ZValue {
	return ERROR("Operator '>' not defined for bool")
}
func (z *ZBool) LessThanEqual(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for bool")
}
func (z *ZBool) GreaterThanEqual(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for bool")
}

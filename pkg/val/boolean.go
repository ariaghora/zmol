package val

// Bool type
type ZBool struct {
	Value bool
}

func (z *ZBool) Type() ZValueType { return "bool" }
func (z *ZBool) Str() string {
	if z.Value {
		return "true"
	}
	return "false"
}
func (z *ZBool) Equals(other ZValue) bool {
	if other.Type() != ZBOOL {
		return false
	}
	return z.Value == other.(*ZBool).Value
}
func (z *ZBool) LessThanEquals(other ZValue) bool {
	if other.Type() != ZBOOL {
		return false
	}
	return !z.Value && other.(*ZBool).Value
}

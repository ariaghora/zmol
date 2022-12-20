package val

// String type
type ZString struct {
	Value string
}

func (z *ZString) Type() ZValueType { return ZSTRING }
func (z *ZString) Str() string      { return z.Value }
func (z *ZString) Equals(other ZValue) bool {
	if other.Type() != ZSTRING {
		return false
	}
	return z.Value == other.(*ZString).Value
}
func (z *ZString) LessThanEquals(other ZValue) bool {
	if other.Type() != ZSTRING {
		return false
	}
	return z.Value <= other.(*ZString).Value
}

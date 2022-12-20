package val

// Error type
type ZError struct {
	Message string
}

func (z *ZError) Type() ZValueType { return "error" }
func (z *ZError) Str() string      { return "ERROR: " + z.Message }
func (z *ZError) Equals(other ZValue) bool {
	if other.Type() != ZERROR {
		return false
	}
	return z.Message == other.(*ZError).Message
}
func (z *ZError) LessThanEquals(other ZValue) bool {
	if other.Type() != ZERROR {
		return false
	}
	return z.Message <= other.(*ZError).Message
}

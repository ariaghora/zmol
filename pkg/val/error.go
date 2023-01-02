package val

// Error type
type ZError struct {
	Message string
}

func ERROR(message string) *ZError {
	return &ZError{Message: message}
}

func OP_ERROR(op string, a ZValue, b ZValue) *ZError {
	return ERROR("Operator '" + op + "' not defined for " + string(a.Type()) + " and " + string(b.Type()))
}

func (z *ZError) Type() ZValueType { return ZERROR }
func (z *ZError) Str() string      { return "ERROR: " + z.Message }
func (z *ZError) Equal(other ZValue) ZValue {
	return ERROR("Cannot compare error with " + string(other.Type()))
}
func (z *ZError) NotEqual(other ZValue) ZValue {
	return ERROR("Cannot compare error with " + string(other.Type()))
}
func (z *ZError) LessThanEqual(other ZValue) ZValue {
	return ERROR("Cannot compare error with " + string(other.Type()))
}
func (z *ZError) GreaterThanEqual(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for error")
}

package val

type ZList struct {
	Elements []ZValue
}

func (zl *ZList) Type() ZValueType {
	return ZLIST
}

func (zl *ZList) Str() string {
	out := "["
	for i, e := range zl.Elements {
		out += e.Str()
		if i < len(zl.Elements)-1 {
			out += ", "
		}
	}
	out += "]"
	return out
}

func (zl *ZList) Literal() string {
	return zl.Str()
}

func (zl *ZList) Equal(other ZValue) ZValue {
	return ERROR("Operator '==' not defined for list")
}

func (zl *ZList) NotEqual(other ZValue) ZValue {
	return ERROR("Operator '!=' not defined for list")
}

func (zl *ZList) LessThan(other ZValue) ZValue {
	return ERROR("Operator '<' not defined for list")
}

func (zl *ZList) GreaterThan(other ZValue) ZValue {
	return ERROR("Operator '>' not defined for list")
}

func (zl *ZList) LessThanEqual(other ZValue) ZValue {
	return ERROR("Operator '<=' not defined for list")
}

func (zl *ZList) GreaterThanEqual(other ZValue) ZValue {
	return ERROR("Operator '>=' not defined for list")
}

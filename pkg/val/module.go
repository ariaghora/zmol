package val

import "fmt"

type ZModule struct {
	ModulePath string
	Env        *Env
}

func (zm *ZModule) Type() ZValueType {
	return ZMODULE
}

func (zm *ZModule) Str() string {
	repr := fmt.Sprintf("<%s path=\"%s\">", zm.Type(), zm.ModulePath)
	return repr
}

func (zm *ZModule) DotAccess(name string) ZValue {
	value, ok := zm.Env.Get(name)
	if !ok {
		return ERROR(fmt.Sprintf("Module '%s' has no attribute '%s'", zm.ModulePath, name))
	}
	return value
}

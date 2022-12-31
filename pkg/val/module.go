package val

import "fmt"

type ZModule struct {
	ModulePath string
	env        *Env
}

func MODULE(modulePath string, env *Env) *ZModule {
	return &ZModule{
		ModulePath: modulePath,
		env:        env,
	}
}

func (zm *ZModule) Type() ZValueType {
	return ZMODULE
}

func (zm *ZModule) Str() string {
	repr := fmt.Sprintf("<%s path=\"%s\">", zm.Type(), zm.ModulePath)
	return repr
}

func (zm *ZModule) DotAccess(name string) ZValue {
	value, ok := zm.env.Get(name)
	if !ok {
		return ERROR(fmt.Sprintf("Module '%s' has no attribute '%s'", zm.ModulePath, name))
	}
	return value
}

func (zm *ZModule) DotAssign(name string, value ZValue) {
	fmt.Println("cannt assign to module")
}

func (zm *ZModule) Env() *Env {
	return zm.env
}

package val

import "fmt"

type ZClass struct {
	Name string
	env  *Env
}

func CLASS(name string, env *Env) *ZClass {
	classDef := &ZClass{
		Name: name,
		env:  env,
	}
	return classDef
}

func (z *ZClass) Type() ZValueType { return ZCLASS }
func (z *ZClass) Str() string      { return fmt.Sprintf("<%s \"%s\">", z.Type(), z.Name) }

func (z *ZClass) DotAccess(name string) ZValue {
	value, ok := z.env.Get(name)
	if !ok {
		return ERROR(fmt.Sprintf("Class '%s' has no attribute '%s'", z.Name, name))
	}
	return value
}

func (z *ZClass) DotAssign(name string, value ZValue) {
	z.env.Set(name, value)
}

func (z *ZClass) Env() *Env {
	return z.env
}

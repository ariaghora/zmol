package val

import "fmt"

type ZObject struct {
	ClassName string
	env       *Env
}

func OBJECT(className string, parentEnv *Env) *ZObject {
	obj := &ZObject{
		ClassName: className,
		env: &Env{
			SymTable:  map[string]ZValue{},
			ParentEnv: parentEnv,
		},
	}
	obj.env.Set("self", obj)
	return obj
}

func (z *ZObject) Type() ZValueType { return ZOBJECT }
func (z *ZObject) Str() string      { return fmt.Sprintf("<%s \"%s\">", z.Type(), z.ClassName) }

func (z *ZObject) DotAccess(name string) ZValue {
	value, ok := z.env.Get(name)
	if !ok {
		return ERROR(fmt.Sprintf("Object '%s' has no attribute '%s'", z.ClassName, name))
	}
	return value
}

func (z *ZObject) DotAssign(name string, value ZValue) {
	if name == "new" {
		fmt.Println("`new` is a reserved attribute name")
		return
	}
	z.env.Set(name, value)
}

func (z *ZObject) Env() *Env {
	return z.env
}

package native

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"

	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/native/goplugin"
	"github.com/ariaghora/zmol/pkg/native/std"
	"github.com/ariaghora/zmol/pkg/val"
)

type NativeFuncRegistry struct {
	zState *eval.ZmolState
}

func NewNativeFuncRegistry(zState *eval.ZmolState) *NativeFuncRegistry {
	return &NativeFuncRegistry{
		zState: zState,
	}
}

func (reg *NativeFuncRegistry) RegisterNativeFunc() {
	// Constructs
	reg.zState.Env.Set("import", &val.ZNativeFunc{Fn: reg.Z_import})
	reg.zState.Env.Set("print", &val.ZNativeFunc{Fn: Z_print})
	reg.zState.Env.Set("println", &val.ZNativeFunc{Fn: Z_println})
	reg.zState.Env.Set("range_list", &val.ZNativeFunc{Fn: Z_range_list})

	// itertools
	reg.zState.Env.Set("append", &val.ZNativeFunc{Fn: Z_append})
	reg.zState.Env.Set("filter", &val.ZNativeFunc{Fn: Z_filter})
	reg.zState.Env.Set("len", &val.ZNativeFunc{Fn: Z_len})
	reg.zState.Env.Set("reduce", &val.ZNativeFunc{Fn: Z_reduce})
	reg.zState.Env.Set("reverse", &val.ZNativeFunc{Fn: Z_reverse})
	reg.zState.Env.Set("zip", &val.ZNativeFunc{Fn: Z_zip})

	// string manipulation
	reg.zState.Env.Set("split", &val.ZNativeFunc{Fn: Z_split})

	// type conversion
	reg.zState.Env.Set("int", &val.ZNativeFunc{Fn: Z_int})
	reg.zState.Env.Set("float", &val.ZNativeFunc{Fn: Z_float})
}

func (reg *NativeFuncRegistry) Z_import(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return &val.ZError{Message: "import takes 1 argument"}
	}

	if args[0].Type() != val.ZSTRING {
		return &val.ZError{Message: "import takes 1 string"}
	}

	// Try import std lib
	switch args[0].(*val.ZString).Value {
	case "goplugin":
		return goplugin.GoPluginModule
	case "io":
		return std.IOModule
	case "math":
		return std.MathModule
	}

	// FIXME: handle !ok
	dir, _ := reg.zState.Env.Get("__moddir__")

	modulePath := args[0].(*val.ZString).Value
	modulePath = path.Join(dir.Str(), modulePath)

	// TODO: Check duplicate import

	content, err := ioutil.ReadFile(modulePath)
	if err != nil {
		msg := fmt.Sprintf("cannot read or load module \"%s\"", modulePath)
		eval.RuntimeErrorf(msg)
		return &val.ZError{Message: msg}
	}

	zState := eval.NewZmolState(nil)
	moduleDir := filepath.Dir(modulePath)
	zState.Env.Set("__moddir__", &val.ZString{Value: moduleDir})
	NewNativeFuncRegistry(zState).RegisterNativeFunc()

	_, err = zState.Eval(string(content))
	if err != nil {
		eval.RuntimeErrorf("cannot eval module " + modulePath)
	}

	return val.MODULE(modulePath, zState.Env)
}

func Z_print(args ...val.ZValue) val.ZValue {
	for _, arg := range args {
		fmt.Print(arg.Str())
	}
	return &val.ZNull{}
}

func Z_println(args ...val.ZValue) val.ZValue {
	for _, arg := range args {
		fmt.Print(arg.Str())
	}
	fmt.Println()
	return &val.ZNull{}
}

func Z_range_list(args ...val.ZValue) val.ZValue {
	if len(args) != 2 {
		return &val.ZError{Message: "range_list takes 2 arguments"}
	}
	start := args[0]
	end := args[1]
	if start.Type() != val.ZINT || end.Type() != val.ZINT {
		return &val.ZError{Message: "range_list takes 2 integers"}
	}

	element := []val.ZValue{}
	for i := start.(*val.ZInt).Value; i < end.(*val.ZInt).Value; i++ {
		element = append(element, &val.ZInt{Value: i})
	}

	list := &val.ZList{
		Elements: element,
	}
	return list
}

func Z_int(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return &val.ZError{Message: "int takes 1 argument"}
	}

	switch args[0].Type() {
	case val.ZINT:
		return args[0]
	case val.ZFLOAT:
		return val.INT(int64(args[0].(*val.ZFloat).Value))
	case val.ZSTRING:
		strval := args[0].(*val.ZString).Value
		res, err := strconv.ParseInt(strval, 10, 64)
		if err != nil {
			eval.RuntimeErrorf("cannot convert string \"" + strval + "\" to int")
		}
		return val.INT(res)
	default:
		eval.RuntimeErrorf("int() takes a number or string as argument")
	}
	return &val.ZNull{}
}

func Z_float(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return &val.ZError{Message: "float takes 1 argument"}
	}

	switch args[0].Type() {
	case val.ZINT:
		return val.FLOAT(float64(args[0].(*val.ZInt).Value))
	case val.ZFLOAT:
		return args[0]
	case val.ZSTRING:
		strval := args[0].(*val.ZString).Value
		res, err := strconv.ParseFloat(strval, 64)
		if err != nil {
			eval.RuntimeErrorf("cannot convert string \"" + strval + "\" to float")
		}
		return val.FLOAT(res)
	default:
		eval.RuntimeErrorf("float() takes a number or string as argument")
	}
	return &val.ZNull{}
}

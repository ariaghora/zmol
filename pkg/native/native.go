package native

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/ariaghora/zmol/pkg/eval"
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

	// IO
	reg.zState.Env.Set("read_string_file", &val.ZNativeFunc{Fn: Z_read_string_file})

	// itertools
	reg.zState.Env.Set("append", &val.ZNativeFunc{Fn: Z_append})
	reg.zState.Env.Set("filter", &val.ZNativeFunc{Fn: Z_filter})
	reg.zState.Env.Set("len", &val.ZNativeFunc{Fn: Z_len})
	reg.zState.Env.Set("reduce", &val.ZNativeFunc{Fn: Z_reduce})
	reg.zState.Env.Set("reverse", &val.ZNativeFunc{Fn: Z_reverse})
	reg.zState.Env.Set("zip", &val.ZNativeFunc{Fn: Z_zip})

	// math
	reg.zState.Env.Set("sqrt", &val.ZNativeFunc{Fn: Z_sqrt})

	// string manipulation
	reg.zState.Env.Set("split", &val.ZNativeFunc{Fn: Z_split})

	// type conversion
	reg.zState.Env.Set("int", &val.ZNativeFunc{Fn: Z_int})
	reg.zState.Env.Set("float", &val.ZNativeFunc{Fn: Z_float})
}

func EnsureFloat(n val.ZValue) (float64, error) {
	if n.Type() == val.ZINT {
		return float64(n.(*val.ZInt).Value), nil
	} else if n.Type() == val.ZFLOAT {
		return n.(*val.ZFloat).Value, nil
	} else {
		return 0, errors.New("not a number")
	}
}

func EnsureInt(n val.ZValue) (int64, error) {
	if n.Type() == val.ZINT {
		return n.(*val.ZInt).Value, nil
	} else if n.Type() == val.ZFLOAT {
		return int64(n.(*val.ZFloat).Value), nil
	}
	return 0, errors.New("not an integer")
}

func (reg *NativeFuncRegistry) Z_import(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return &val.ZError{Message: "import takes 1 argument"}
	}

	if args[0].Type() != val.ZSTRING {
		return &val.ZError{Message: "import takes 1 string"}
	}

	modulePath := args[0].(*val.ZString).Value
	content, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return &val.ZError{Message: "cannot read or load module \"" + modulePath + "\""}
	}

	zState := eval.NewZmolState(nil)
	NewNativeFuncRegistry(zState).RegisterNativeFunc()
	zState.Eval(string(content))

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

func Z_read_string_file(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return &val.ZError{Message: "read_string_file takes 1 argument"}
	}

	if args[0].Type() != val.ZSTRING {
		return &val.ZError{Message: "read_string_file takes 1 string"}
	}

	filePath := args[0].(*val.ZString).Value
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		eval.RuntimeErrorf("cannot read file " + filePath)
	}
	return val.STRING(string(content))
}

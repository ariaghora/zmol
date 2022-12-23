package native

import (
	"errors"
	"fmt"

	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/val"
)

func RegisterNativeFunc(zState *eval.ZmolState) {
	// IO
	zState.Env.Set("print", &val.ZNativeFunc{Fn: Z_print})
	zState.Env.Set("println", &val.ZNativeFunc{Fn: Z_println})
	zState.Env.Set("range_list", &val.ZNativeFunc{Fn: Z_range_list})

	// itertools
	zState.Env.Set("append", &val.ZNativeFunc{Fn: Z_append})
	zState.Env.Set("filter", &val.ZNativeFunc{Fn: Z_filter})
	zState.Env.Set("len", &val.ZNativeFunc{Fn: Z_len})
	zState.Env.Set("reduce", &val.ZNativeFunc{Fn: Z_reduce})
	zState.Env.Set("reverse", &val.ZNativeFunc{Fn: Z_reverse})
	zState.Env.Set("zip", &val.ZNativeFunc{Fn: Z_zip})

	// math
	zState.Env.Set("sqrt", &val.ZNativeFunc{Fn: Z_sqrt})

	// type conversion
	zState.Env.Set("int", &val.ZNativeFunc{Fn: Z_int})
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

func Z_print(args ...val.ZValue) val.ZValue {
	for _, arg := range args {
		fmt.Print(arg.Str())
	}
	return &val.ZNull{}
}

func Z_println(args ...val.ZValue) val.ZValue {
	for _, arg := range args {
		fmt.Println(arg.Str())
	}
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

	n, err := EnsureInt(args[0])
	if err != nil {
		return &val.ZError{Message: "int takes an integer as argument"}
	}

	return val.INT(n)
}

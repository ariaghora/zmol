package native

import (
	"fmt"

	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/val"
)

func RegisterNativeFunc(zState *eval.ZmolState) {
	zState.Env.Set("print", &val.ZNativeFunc{Fn: Z_print})
	zState.Env.Set("range_list", &val.ZNativeFunc{Fn: Z_range_list})
	zState.Env.Set("filter", &val.ZNativeFunc{Fn: Z_filter})
}

func Z_print(args ...val.ZValue) val.ZValue {
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

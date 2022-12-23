package native

import (
	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/val"
)

func Z_filter(args ...val.ZValue) val.ZValue {
	if len(args) != 2 {
		return &val.ZError{Message: "filter takes 2 arguments"}
	}
	list := args[0]
	fn := args[1]
	if list.Type() != val.ZLIST || fn.Type() != val.ZFUNCTION {
		return &val.ZError{Message: "filter takes a list and a function"}
	}

	// check if number of arguments in function is 1
	if len(fn.(*val.ZFunction).Params) != 1 {
		return &val.ZError{Message: "filter takes a function with 1 argument"}
	}

	actualFn := fn.(*val.ZFunction)
	zState := eval.NewZmolState(actualFn.Env)

	element := []val.ZValue{}

	for _, e := range list.(*val.ZList).Elements {
		actualFn.Env.Set(actualFn.Params[0].Value, e)
		isTrue := zState.EvalProgram(actualFn.Body)

		// check if return value is boolean
		if isTrue.Type() != val.ZBOOL {
			return &val.ZError{Message: "filter takes a function that returns a boolean"}
		}
		if isTrue.(*val.ZBool).Value {
			element = append(element, e)
		}
	}

	list = &val.ZList{
		Elements: element,
	}
	return list
}

func Z_len(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return &val.ZError{Message: "len takes 1 argument"}
	}
	list := args[0]
	if list.Type() != val.ZLIST || list.Type() != val.ZSTRING {
		return &val.ZError{Message: "len takes a list"}
	}

	if list.Type() == val.ZLIST {
		return &val.ZInt{Value: int64(len(list.(*val.ZList).Elements))}
	} else {
		return &val.ZInt{Value: int64(len(list.(*val.ZString).Value))}
	}
}

func Z_reduce(args ...val.ZValue) val.ZValue {
	if len(args) != 3 {
		return &val.ZError{Message: "reduce takes 3 arguments"}
	}
	list := args[0]
	fn := args[1]
	initial := args[2]
	if list.Type() != val.ZLIST || fn.Type() != val.ZFUNCTION {
		return &val.ZError{Message: "reduce takes a list and a function"}
	}

	// check if number of arguments in function is 2
	if len(fn.(*val.ZFunction).Params) != 2 {
		return &val.ZError{Message: "reduce takes a function with 2 arguments"}
	}

	actualFn := fn.(*val.ZFunction)
	zState := eval.NewZmolState(actualFn.Env)

	element := initial

	for _, e := range list.(*val.ZList).Elements {
		actualFn.Env.Set(actualFn.Params[0].Value, element)
		actualFn.Env.Set(actualFn.Params[1].Value, e)
		element = zState.EvalProgram(actualFn.Body)
	}

	return element
}

func Z_append(args ...val.ZValue) val.ZValue {
	if len(args) != 2 {
		return &val.ZError{Message: "append takes 2 arguments"}
	}
	list := args[0]
	element := args[1]
	if list.Type() != val.ZLIST {
		return &val.ZError{Message: "append takes a list as first argument"}
	}

	list.(*val.ZList).Elements = append(list.(*val.ZList).Elements, element)
	return list
}

func Z_reverse(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return &val.ZError{Message: "reverse takes 1 argument"}
	}
	list := args[0]
	if list.Type() != val.ZLIST {
		return &val.ZError{Message: "reverse takes a list"}
	}

	element := []val.ZValue{}
	for i := len(list.(*val.ZList).Elements) - 1; i >= 0; i-- {
		element = append(element, list.(*val.ZList).Elements[i])
	}

	list = &val.ZList{
		Elements: element,
	}
	return list
}

func Z_zip(args ...val.ZValue) val.ZValue {
	if len(args) != 2 {
		return &val.ZError{Message: "zip takes 2 arguments"}
	}
	list1 := args[0]
	list2 := args[1]
	if list1.Type() != val.ZLIST || list2.Type() != val.ZLIST {
		eval.RuntimeErrorf("zip takes 2 lists")
	}

	if len(list1.(*val.ZList).Elements) != len(list2.(*val.ZList).Elements) {
		return eval.RuntimeErrorf("zip takes 2 lists with same length")
	}

	element := []val.ZValue{}
	for i := 0; i < len(list1.(*val.ZList).Elements); i++ {
		element = append(element, &val.ZList{
			Elements: []val.ZValue{
				list1.(*val.ZList).Elements[i],
				list2.(*val.ZList).Elements[i],
			},
		})
	}

	list := &val.ZList{
		Elements: element,
	}
	return list
}

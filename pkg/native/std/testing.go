package std

import (
	"fmt"

	"github.com/ariaghora/zmol/pkg/val"
)

var TestingModule = val.MODULE(
	"testing",
	&val.Env{
		SymTable: map[string]val.ZValue{
			"assert_true":  &val.ZNativeFunc{Fn: Z_testing_assert_true},
			"assert_equal": &val.ZNativeFunc{Fn: Z_testing_assert_equal},
		},
	},
)

func Z_testing_assert_true(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return val.ERROR("assert_true() takes exactly 1 argument")
	}

	if args[0].Type() != val.ZBOOL {
		return val.ERROR("assert_true() takes a boolean as argument")
	}

	if !args[0].(*val.ZBool).Value {
		return val.ERROR("assertion failed, expected true but got false")
	}

	return val.NULL()
}

func Z_testing_assert_equal(args ...val.ZValue) val.ZValue {
	if len(args) != 2 {
		return val.ERROR("assert_equal() takes exactly 2 arguments")
	}

	// check if both arguments are ZComparable
	left, leftOk := args[0].(val.ZComparable)
	right, rightOk := args[1].(val.ZComparable)
	if !leftOk || !rightOk {
		return val.ERROR("assert_equal() takes comparable arguments")
	}

	if !left.Equals(right).(*val.ZBool).Value {
		return val.ERROR(fmt.Sprintf("assertion failed, expected %v but got %v", left.Str(), right.Str()))
	}

	return val.NULL()
}

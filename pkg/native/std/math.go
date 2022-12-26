package std

import (
	"errors"
	"math"

	"github.com/ariaghora/zmol/pkg/val"
)

// This math module is primarily to wrap Go's math optimized functions.
var MathModule = val.MODULE(
	"math",
	&val.Env{
		SymTable: map[string]val.ZValue{
			// Constants
			"PI": val.FLOAT(math.Pi),
			"E":  val.FLOAT(math.E),
			// Unary functions
			"abs":   Z_floatUfunc("abs", math.Abs),
			"acos":  Z_floatUfunc("acos", math.Acos),
			"asin":  Z_floatUfunc("asin", math.Asin),
			"atan":  Z_floatUfunc("atan", math.Atan),
			"cos":   Z_floatUfunc("cos", math.Cos),
			"exp":   Z_floatUfunc("exp", math.Exp),
			"log":   Z_floatUfunc("log", math.Log),
			"log2":  Z_floatUfunc("log2", math.Log2),
			"log10": Z_floatUfunc("log10", math.Log10),
			"sin":   Z_floatUfunc("sin", math.Sin),
			"sqrt":  Z_floatUfunc("sqrt", math.Sqrt),
			"tan":   Z_floatUfunc("tan", math.Tan),
			// More functions
			// ...
		},
	},
)

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

// Creates a function that takes a single float argument and returns a float.
// Typically used for Go's math unary functions like sqrt, exp, etc.
func Z_floatUfunc(name string, fn func(float64) float64) val.ZValue {
	return &val.ZNativeFunc{
		Fn: func(args ...val.ZValue) val.ZValue {
			if len(args) != 1 {
				return val.ERROR(name + "() takes exactly 1 argument")
			}

			n, err := EnsureFloat(args[0])
			if err != nil {
				return val.ERROR(name + "() takes a number as argument")
			}

			n = fn(n)
			return val.FLOAT(n)
		},
	}
}

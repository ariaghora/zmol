package native

import (
	"math"

	"github.com/ariaghora/zmol/pkg/val"
)

func Z_sqrt(args ...val.ZValue) val.ZValue {
	if len(args) != 1 {
		return val.ERROR("sqrt() takes exactly 1 argument")
	}

	n, err := EnsureFloat(args[0])
	if err != nil {
		return val.ERROR("sqrt() takes a number as argument")
	}

	n = math.Sqrt(n)
	return val.FLOAT(n)
}

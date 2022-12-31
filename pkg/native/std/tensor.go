package std

import (
	"github.com/ariaghora/zmol/pkg/val"
	"gorgonia.org/tensor"
)

var TensorModule = val.MODULE(
	"tensor",
	&val.Env{
		SymTable: map[string]val.ZValue{
			"zeros": &val.ZNativeFunc{Fn: Z_tensor_zeros},
		},
	},
)

func Z_tensor_zeros(args ...val.ZValue) val.ZValue {
	shape := args

	// ensure all shape values are integers
	for _, v := range shape {
		if v.Type() != val.ZINT {
			return &val.ZError{Message: "shape values must be integers"}
		}
	}

	// convert shape to []int
	shapeInt := make([]int, len(shape))
	for i, v := range shape {
		shapeInt[i] = int(v.(*val.ZInt).Value)
	}

	tensor := tensor.New(
		tensor.Of(tensor.Float32), tensor.WithShape(shapeInt...),
	)

	return &val.ZTensor{
		Data: tensor,
	}
}

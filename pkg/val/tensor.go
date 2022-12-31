package val

import (
	"fmt"

	"gorgonia.org/tensor"
)

type ZTensor struct {
	Data tensor.Tensor
}

func (z *ZTensor) Type() ZValueType { return ZTENSOR }
func (z *ZTensor) Str() string      { return fmt.Sprintf("%v", z.Data) }

package native

import (
	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/val"
)

func RegisterNativeFunc(zState *eval.ZmolState) {
	zState.Env.Set("print", &val.ZNativeFunc{Fn: Z_Print})
}

func Z_Print(args ...val.ZValue) val.ZValue {
	for _, arg := range args {
		print(arg.Str())
	}
	return &val.ZNull{}
}

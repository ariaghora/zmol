package std

import (
	"io/ioutil"

	"github.com/ariaghora/zmol/pkg/eval"
	"github.com/ariaghora/zmol/pkg/val"
)

var IOModule = val.MODULE(
	"io",
	&val.Env{
		SymTable: map[string]val.ZValue{
			"read_string_file": &val.ZNativeFunc{Fn: Z_read_string_file},
		},
	},
)

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

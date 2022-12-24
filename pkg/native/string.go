package native

import (
	"strings"

	"github.com/ariaghora/zmol/pkg/val"
)

func Z_split(args ...val.ZValue) val.ZValue {
	if len(args) != 2 {
		return val.ERROR("split() takes exactly 2 arguments")
	}

	if args[0].Type() != val.ZSTRING && args[1].Type() != val.ZSTRING {
		return val.ERROR("split() takes a string as argument")
	}

	s := args[0].(*val.ZString).Value
	sep := args[1].(*val.ZString).Value

	elements := strings.Split(s, sep)
	list := []val.ZValue{}
	for _, e := range elements {
		list = append(list, val.STRING(e))
	}

	return &val.ZList{
		Elements: list,
	}

}
